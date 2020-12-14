# WALLE
4swap mtg gateway（4swap mtg 程序化交易支付网关）

### MTG Memo

与给普通 4swap 转账下单不同，4swap mtg 程序下单是转账给一个多签地址，因为 memo 暴露在主网的 Tx 里，为了保护用户隐私需要对 memo 进行加密处理。

mtg memo 每个部分的第一个 byte 的值代表这部分内容的 byte 长度（例如 uuid 占 17 位 bytes，第一位 byte 值是 16，后面 16 位是 uuid 的内容），各部分按约定顺序拼接得到未加密的 raw memo。

生成一对 ed25519 公私钥对，使用生成的私钥和 4swap mtg 公开的签名公钥合成一对 aes 加密需要的 key 和 iv，加密上述 raw memo，将加密后的
memo 拼接在生成的公钥后面，最后使用 base64 编码处理得到最后的 memo。

memo 生成参考 [4swap-sdk-go](https://github.com/fox-one/4swap-sdk-go/blob/master/mtg/action.go#L71)

### 多签转账下单

通过 ```/api/info``` 拿到多签信息，然后通过 [mixin api POST /transactions](https://github.com/fox-one/mixin-sdk-go/blob/master/transaction_raw.go#L42) 转账。

```json5
// /api/info
{
  "ts": 1607936726300,
  "data": {
    "members": [
      "9656eacd-2fa7-4e7b-b0eb-c475c9964f78",
      "ab14736f-e595-4e65-9879-871819d390f5",
      "b856deb3-e92f-4c19-9733-ec43526f95ce",
      "229fc7ac-9d09-4a6a-af5a-78f7439dce76",
      "84a4db41-4992-4d35-aac7-987f965f0302"
    ],
    "public_key": "WE2b3mzyi23SiEKEiiHy6+72LVUG9gDSEJ0d1jU+yC0=", // 用于生成 memo 加密需要的 aes key & iv 的公钥
    "threshold": 4
  }
}
```

### gateway 的作用

帮助 4swap 程序化交易程序无缝迁移到 4swap mtg 交易

1. 给 gateway 机器人发送文本消息 ```broker``` 申请自己专属的 broker 钱包
2. 将之前直接发给 4swap 机器人的转账，改成发给上面申请的 broker，broker 收到的转账之后会立即将 memo 进行编码加密之后转发到 mtg 多签地址
3. 转账给 broker 的 trace id 依然是订单 id，可以继续通过 ```/api/orders/id``` api 查询订单信息
4. 买到或者退款的币将直接从多签钱包转到下单的钱包，无需通过 gateway 中转。
