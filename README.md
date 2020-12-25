# 4swap mtg gateway


4swap mtg 支付网关的作用是将普通版下单的转账升级并转发到多签版本。

### 使用方式

* 给 gateway 机器人(mixin id 7000103594)发送文本消息 ```broker``` 申请自己专属的 broker 钱包
* 将之前直接发给 4swap 普通版的转账，改成发给上面申请的 broker，broker 收到的转账之后会立即将 memo 进行编码加密之后转发到 mtg 多签地址
* 转账给 broker 的 trace id 依然是订单 id，可以继续通过 ```/api/orders/id``` api 查询订单信息
* 买到或者退款的币将直接从多签钱包转到下单的钱包，无需通过 gateway 中转。

### 补充

* 同一用户可以多次申请，申请到的 broker 都是有效的。
* 多签版本的 4swap api host 是 https://f1-mtgswap-api.firesbox.com
* 普通版本的 4swap api host 是 https://f1-uniswap-api.firesbox.com
* 了解多签版本 memo 生成方式，请参考 [4swap-sdk-go](https://github.com/fox-one/4swap-sdk-go/blob/master/mtg/action.go)
