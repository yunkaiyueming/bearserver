特殊牌：

改颜色,不改头像
禁卡牌
交换牌
加2禁卡牌
王牌=变色牌

颜色相同，头像相同

4中颜色


====================================
模块：
	注册，登录，
	在线用户
	创建房间
	发牌规则实现
	推送逻辑：推送给其他3个人；长时间没法，服务器自动发；

	输赢规则判断：
	特殊：意外掉线


===============进度=======================
规范
玩法确认&整套卡牌数 40张 10*4种																																																	
分工



玩法：默认5张牌，开一张，默认第一个人出牌，按进房间顺序出牌，

牌模块

===服务器发送规则==
这是个例子，游戏逻辑必须用dispatch
{"Dispatch":{"params":"haha","cmd":"hello","rnum":6,"ts":1499945446}}
这个是返回
{"Response":{"Uid":0,"Cmd":"hello",Ret":0,"Data":{"ID":1,"Name":"Reds"},"Rnum":6}}

=== 接口定义 ===
登录接口:
{"UserLoginInfo":{"params":{},"cmd":"login","name":"sq","rnum":6,"ts":1499945446}}
根据uid推送固定的消息:
{"Dispatch":{"params":{"msg":"hehe"},"uid":1,"cmd":"pushMsg","rnum":6,"ts":1499945446}}

注册接口:
{"RegisterUserInfo":{"name":"xiaocai","pwd":"123"}}

发牌接口:
发牌：
{"Dispatch":{"params":{"mtype":"1","card":"401"},"uid":1,"cmd":"playCard","rnum":6,"ts":1499945446}}
摸牌:
{"Dispatch":{"params":{"mtype":"2"},"uid":1,"cmd":"playCard","rnum":6,"ts":1499945446}}





