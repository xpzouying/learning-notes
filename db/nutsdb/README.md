# NutsDB

代码使用 `tag=v0.1.0` 版本。

**db 两种方式**

- HintAndRAMIdxMode - 默认，Entry 整个都保存在内存中。
- HintAndMemoryMapIdxMode - 内存中就保存索引，真正的数据需要从数据文件中获取。


## <a id="contents"> Contents </a>

- [读写操作](#simple-put-get)


<br /> <hr />



## <a id="simple-put-get"> 读写操作 </a>

<details>

<summary> 代码示例 </summary>

```go
package main

import (
	"log"

	"github.com/xujiajun/nutsdb"
)

func main() {
	opt := nutsdb.DefaultOptions
	opt.Dir = "/tmp/nutsdb" // 这边数据库会自动创建这个目录文件
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// to do
	if err := db.Update(func(tx *nutsdb.Tx) error {
		if err := tx.Put("bucket1", []byte("foo"), []byte("bar"), 0); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get("bucket1", []byte("foo"))
		if err != nil {
			return err
		}

		log.Printf("entry: %+v", entry)
		return nil
	}); err != nil {
		panic(err)
	}

}
```

</details>


`Update` 和 `View` 方法，是为了保障事务，暂时先跳过，先关注读写操作流程以及对应的流程。


<br /> <hr />

### <a id="entry-object">Entry 对象</a>

一个操作对象，在 NutsDB 中叫 Entry，其数据结构如下：

![](./assets/nutsdb_record.jpg)

这里也分为 Header 和 Data 两块区域：

1. Header - 40 Bytes（固定长度）。从 crc 到 txId。
2. Data - bucket、key、value。


<br /> <hr />

### 底层实现 - 写

```go
func (tx *Tx) put(bucket string, key, value []byte, ttl uint32, flag uint16, timestamp uint64) error {
    // 删除无关代码

	tx.pendingWrites = append(tx.pendingWrites, &Entry{
		Key:   key,
		Value: value,
		Meta: &MetaData{
			keySize:    uint32(len(key)),
			valueSize:  uint32(len(value)),
			timestamp:  timestamp,
			Flag:       flag,
			TTL:        ttl,
			bucket:     []byte(bucket),
			bucketSize: uint32(len(bucket)),
			status:     UnCommitted,
			txId:       tx.id,
		},
	})

	return nil
}
```

1. `put` 方法是对该事务的 Entry 列表，增加一个 Entry 对象而已。再未完成提交之前，都还在内存中。

2. 在 `Tx.Commit` 时，才会真正的写盘。Commit 在 Update 方法中，在没有错误的情况下，会自动调用完成。
   事务的分析放到后面介绍。


<br /> <hr />

### 底层实现 - 读


**HintAndRAMIdxMode 模式**


默认使用的 db 模式是：`HintAndRAMIdxMode`。暂时删除其他无关代码。

```go
// Get retrieves the value for a key in the bucket.
// The returned value is only valid for the life of the transaction.
func (tx *Tx) Get(bucket string, key []byte) (e *Entry, err error) {

	// 先找到索引， idx 索引。
	if idx, ok := tx.db.HintIdx[bucket]; ok {

		r, err := idx.Find(key)
		if err != nil {
			return nil, err
		}

		if r.H.meta.Flag == DataDeleteFlag || r.isExpired() {
			return nil, ErrNotFoundKey
		}

		if idxMode == HintAndRAMIdxMode {
			return r.E, nil
		}
	}

	return nil, errors.New("not found bucket:" + bucket + ",key:" + string(key))
}
```

1. 通过 bucket 名字，从 `hash` 找到对应 bucket 的索引，索引使用 `B+` 树实现。

2. 然后通过索引找到对应的 Entry 。因为是使用 `HintAndRAMIdxMode` 模式，所以找到的 Record。其中 Record 包含对应的 Metadata 和 Entry，通过 Metadata 判断是否有效，若未有效，则返回里面的 Entry 即可。

<br />

**B+ 树的数据**

所有的数据都保存在 B+ 树的叶子结点中，叶子结点上面的元素叫做：`Record`，其中的结构为：

```go
// Record records entry and hint.
Record struct {
	H *Hint
	E *Entry
}
```

1. Hint - 数据元素的 MetaData 数据。包含，是否被删除、过期时间等等。
2. Entry - 真实的数据。


> TODO(zy): 在 Hint 和 Entry 中都包含 MetaData，不知道为什么需要相同的数据？


<br />

**HintAndMemoryMapIdxMode 模式**

如果开启 `HintAndMemoryMapIdxMode` 模式，则内存中只是包含索引，真实的数据都保存在文件中。

```go
// Get retrieves the value for a key in the bucket.
// The returned value is only valid for the life of the transaction.
func (tx *Tx) Get(bucket string, key []byte) (e *Entry, err error) {
	if idx, ok := tx.db.HintIdx[bucket]; ok {  // B+ 树找到对应的索引

		r, err := idx.Find(key)  // 
		// 省略部分代码

		if idxMode == HintAndMemoryMapIdxMode {
			path := tx.db.getDataPath(r.H.fileId)
			df, err := NewDataFile(path, tx.db.opt.SegmentSize)
			if err != nil {
				return nil, err
			}

			item, err := df.ReadAt(int(r.H.dataPos))
			if err != nil {
				return nil, fmt.Errorf("read err. pos %d, key %s, err %s", r.H.dataPos, string(key), err)
			}

			if err := df.m.Unmap(); err != nil {
				return nil, err
			}

			return item, nil
		}
	}

	return nil, errors.New("not found bucket:" + bucket + ",key:" + string(key))
}
```

先找到对应的 B+ 实现的索引，这一步与前面的一致。

由于内存中的 B+ 树，只包含对应的数据索引，真正的数据需要去文件中读取。所以这里会通过 Record 中的 Hint 数据获取对应 Entry 所在的文件和偏移位，读出真正的数据。

<br />

**如何从磁盘文件中加载数据**

1. 从 B+ 树索引获取文件和数据在文件中的偏移位（Offset），即可定位到 Entry 在文件的开始位置。
2. [Entry 结构体](#entry-object) 分为 Header 和 Data 两块。首先加载 Header （固定的 40 字节），然后根据 Header 中的数据 `bucketSize`、 `keySize`、`valueSize`，分别加载出对应的数据。

> TODO(zy): 在这里使用了 mmap 的文件映射，进行加速。需要研究一下 mmap 的使用。


<br /> <hr />


## 事务提交

在前面分析 `put` 方法，进行写的时候，实际上只是完成一部分，只是暂时将数据写到了内存中，并没有完成一次完整的事务操作。在此时，如果有其他的写成读取对应数据时，实际上是获取不到对应的值。只有让事务完成 `Commit` 操作后，才算是真正的结束写操作的全部流程。

完成事务调教大概的流程为：

1. 先检查是否有新的数据待写入，即前面提到的： `pendingWrites`。
2. 检查写入的文件是否有足够的空间保存。如果没有，则需要 `rotateActiveFile`。
3. 将新的数据写入到文件中。
4. 更新索引中的 Record 中的 Hint 数据。

**源码分析**

假设文件有足够的空间保存我们待写入的数据，先看看如何将数据写入到文件中，并且更新索引中的数据。

默认的文件大小是：**64 MB**。

```go
// 删除部分代码
func (tx *Tx) Commit() error {
	var e *Entry

	// ...

	for i := 0; i < writesLen; i++ {
		entry := tx.pendingWrites[i]
		
		// 如果是最后一个 待写入 的数据
		if i == writesLen-1 {
			entry.Meta.status = Committed
		}

		// 当前文件的写入的 offset
		off := tx.db.ActiveFile.writeOff
		// 将 entry 的数据，写入到文件的位置上。（写磁盘文件完成！！！）
		if _, err := tx.db.ActiveFile.WriteAt(entry.Encode(), off); err != nil {
			return err
		}

		// 更新磁盘文件的状态
		tx.db.ActiveFile.ActualSize += entrySize
		tx.db.ActiveFile.writeOff += entrySize

		// 如果是内存级别的，则每一个 entry 都会标记为 Committed
		if tx.db.opt.EntryIdxMode == HintAndRAMIdxMode {
			entry.Meta.status = Committed
			e = entry
		} else {
			// 如何是磁盘持久化级别的话，那么这里没有特殊标记。
			// 从函数一进入的代码看，只有最后一个 entry 才会标记为 Committed。
			e = nil
		}

		countFlag := CountFlagEnabled
		// ...

		// 每一个 bucket 是有一个 B+ 树索引
		bucket := string(entry.Meta.bucket)
		if _, ok := tx.db.HintIdx[bucket]; !ok {
			tx.db.HintIdx[bucket] = NewTree()
		}

		// 在 B+ 树索引更新对应的索引属性。
		_ = tx.db.HintIdx[bucket].Insert(entry.Key, e, &Hint{
			fileId:  tx.db.ActiveFile.fileId,
			key:     entry.Key,
			meta:    entry.Meta,
			dataPos: uint64(off),
		}, countFlag)

		tx.db.KeyCount++  // db 中的 key 数量统计
	}

	tx.unlock()

	tx.db = nil

	return nil
}
```

在这里有几个问题（注意：基于 tag=v0.1.0 版本）

**问题 1**

从代码看起来，无论是哪一种索引模式，内存级别的 B+ 树都保存得有完整的 Entry Value。

在这里就比较奇怪，如果是纯内存模型，那么内存中保存完整的 Entry（Key+Value）没有任何问题。
但是，如果是内存+文件混合的模式，那么内存中为什么还需要保存完整的 Entry 呢，因为在读取 Entry 数据的时候，仍然是读取文件中的数据。

> 这里也跟作者聊了一下。答复说：如果是纯内存模型，那么内存里面保存的是完整的 Entry （Key+Value），
> 如果是 `内存索引 + 文件数据` 的混合模式，则内存中应该只是保存索引，数据应该都在文件中。
> 
> 但是从目前的写行为中看，并非如此。看看是不是后面的版本有优化。


**问题 2**

为什么纯内存模式，是每一个 Entry 都会标记为 Committed。

而内存+文件模式，是最后一个 Entry 标记为 Committed？

> 跟作者 [@xujiajun](github.com/xujiajun) 聊了一下，他说按理说都应该只是最后一个提交成功了，才算是一次事务所有的 Entry 都算是 Committed。
>
> 在代码中，没有发现相关的处理。不知道是 v0.1.0 版本的问题，还是漏看的原因。

**问题 3**

因为每一个 bucket 都是一个 B+ 树，如果是并发地更新同一个 bucket 中的数据的话，那么不会产生并发的问题吗？

> 经过测试，发现并发并不会出现任何数据问题。
> 继续搜索代码，发现在开始 Tx 时，`Begin` 方法中，会调用 `tx.lock()`，在该函数中，会对 db 进行锁的保护。
> 如果是 `Update` 的事务，则会开启 `Lock`。如果是 `View` 的事物，则会开启 `RLock`。


<br /> <hr />


## <a id="refs"> Refs </a>

- [NutsDB 说明文档](https://github.com/nutsdb/nutsdb/blob/master/README-CN.md)


## todos

- [ ] 事务提交是如何实现的。
- [ ] TTL 是如何实现的。
- [ ] 写日志，是如何工作的。Append Log 或者 WAL。
- [ ] 具体研究一下 [隔离级别](https://github.com/nutsdb/nutsdb/blob/master/README-CN.md#%E9%9A%94%E7%A6%BB%E7%BA%A7%E5%88%AB%E4%BD%8E%E5%88%B0%E9%AB%98)
- [ ] B+ 树 vs LSM 树 的研究。