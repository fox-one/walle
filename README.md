# WALLE (mixin id 7000103594)
4swap mtg gateway（4swap mtg 程序化交易支付网关）。只需要申请一个 broker，将之前发给 4swap 机器人的转账发给这个 broker，memo 等都不需要改变，
即可完成 4swap mtg 下单。 

### MTG Memo

与给普通 4swap 转账下单不同，4swap mtg 程序下单是转账给一个多签地址，因为 memo 暴露在主网的 Tx 里，为了保护用户隐私需要对 memo 进行加密处理。

mtg memo 每个部分的第一个 byte 的值代表这部分内容的 byte 长度（例如 uuid 占 17 位 bytes，第一位 byte 值是 16，后面 16 位是 uuid 的内容），各部分按约定顺序拼接得到未加密的 raw memo。

生成一对 ed25519 公私钥对，使用生成的私钥和 4swap mtg 公开的签名公钥合成一对 aes 加密需要的 key 和 iv，加密上述 raw memo，将加密后的
memo 拼接在生成的公钥后面，最后使用 base64 编码处理得到最后的 memo。

memo 生成参考 [4swap-sdk-go](https://github.com/fox-one/4swap-sdk-go/blob/master/mtg/action.go#L71)

### 多签转账下单

通过 ```https://f1-mtgswap-api.firesbox.com/api/info``` 拿到多签信息，然后通过 mixin api [POST /transactions](https://github.com/fox-one/mixin-sdk-go/blob/master/transaction_raw.go#L42) 转账。

```json5
// /api/info
{
  "ts": 1608836959157,
  "data": {
    // 节点成员 client id
    "members": [
      "a753e0eb-3010-4c4a-a7b2-a7bda4063f62",
      "099627f8-4031-42e3-a846-006ee598c56e",
      "aefbfd62-727d-4424-89db-ae41f75d2e04",
      "d68ca71f-0e2c-458a-bb9c-1d6c2eed2497",
      "e4bc0740-f8fe-418c-ae1b-32d9926f5863"
    ],
    "public_key": "dt351xp3KjNlVCMqBYUeUSF45upCEiReSZAqcjcP/Lc=", // 用于生成 memo 加密需要的 aes key & iv 的公钥
    "threshold": 3 // 多签组资产转出需要的签名数
  }
}

```

### MTG Gateway

mtg 支付网关的作用是帮助 4swap 程序化交易程序无缝迁移到 4swap mtg 交易

1. 给 gateway 机器人发送文本消息 ```broker``` 申请自己专属的 broker 钱包
2. 将之前直接发给 4swap 机器人的转账，改成发给上面申请的 broker，broker 收到的转账之后会立即将 memo 进行编码加密之后转发到 mtg 多签地址
3. 转账给 broker 的 trace id 依然是订单 id，可以继续通过 ```/api/orders/id``` api 查询订单信息
4. 买到或者退款的币将直接从多签钱包转到下单的钱包，无需通过 gateway 中转。
