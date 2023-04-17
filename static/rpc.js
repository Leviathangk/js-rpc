function GkRpc(url) {
    if (!url) {
        throw new Error("无建立连接的 Url!")
    }

    this.url = url;
    this.functions = {};
    this.uuid = null;
    this.starting = false;
    this.connect();

    this.TypeOpen = 0;   // 打开信号
    this.TypeCreate = 1;   // 创建信号
    this.TypeShow = 2; // 查看函数名信号
    this.TypeDelete = 3; // 删除信号
    this.TypeRun = 4;   // 执行信号
    this.TypeEval = 5;  // eval 信号
}

// connect 创建连接
GkRpc.prototype.connect = function () {
    let _this = this;
    console.log("正在连接：", this.url);
    this.socket = new WebSocket(this.url);

    // onOpen 创建连接时的函数
    this.socket.onopen = function () {
        console.log("连接成功：", _this.url);
    };

    // onmessage 接收信息时的函数
    this.socket.onmessage = function (recv) {
        console.log("接收到数据：", recv.data);
        let recvJson = JSON.parse(recv.data);
        let msgType = recvJson["type"];

        // 判断消息类型
        switch (msgType) {
            case _this.TypeOpen:
                _this.uuid = recv["uuid"];
                _this.sendJson({ "type": _this.TypeOpen, "msg": { "domain": location.hostname } });
                _this.starting === false;
                break
            case _this.TypeCreate:
                _this.createRpc(recvJson);
                break
            case _this.TypeShow:
                _this.show(recvJson);
                break
            case _this.TypeDelete:
                _this.delete(recvJson);
                break
            case _this.TypeRun:
                _this.run(recvJson);
                break
            case _this.TypeEval:
                _this.runEval(recvJson);
                break
            default:
                console.log(recvJson);
        }
    }

    // onerror 连接出错时的函数
    this.socket.onerror = function (e) {
        console.log("连接出错，将在 8s 后重试连接：", e)
        setTimeout(function () {
            if (_this.starting === false) {
                _this.starting === true
                console.log("error 尝试中...")
                _this.connect()
            } else {
                console.log("error 放弃...")
            }
        }, 8 * 1000)
    }

    // onclose 连接关闭时的函数
    this.socket.onclose = function () {
        console.log("连接被动关闭，将在 10s 后重试连接")
        setTimeout(function () {
            if (_this.starting === false) {
                _this.starting === true
                console.log("close 尝试中...")
                _this.connect()
            } else {
                console.log("close 放弃...")
            }
        }, 10 * 1000)
    }
}

// createRpc 创建 rpc
GkRpc.prototype.createRpc = function (recvJson) {
    let _this = this;
    let funcName = recvJson["funcName"];
    const eventId = recvJson["eventId"];
    const funcBody = new Function("sendResult", "recvJson", decodeURIComponent(recvJson["funcBody"]));

    let exists = funcName in _this.functions;

    try {
        _this.functions[funcName] = funcBody;
        if (!exists) {
            _this.sendJson({ "type": _this.TypeCreate, "msg": { "success": true, "msg": "注入成功：" + funcName }, "eventId": eventId })
        } else {
            _this.sendJson({ "type": _this.TypeCreate, "msg": { "success": true, "msg": "覆写成功：" + funcName }, "eventId": eventId })
        }
    } catch (e) {
        _this.sendJson({ "type": _this.TypeCreate, "msg": { "success": false, "msg": "注入错误：" + e.message }, "eventId": eventId })
    }
}

// 查看函数名
GkRpc.prototype.show = function (recvJson) {
    const eventId = recvJson["eventId"];

    let nameArr = [];
    for (let n in this.functions) {
        nameArr.push(n)
    }

    this.sendJson({ "type": this.TypeShow, "msg": { "success": true, "msg": { "functions": nameArr, "total": nameArr.length } }, "eventId": eventId })
}

// 删除指定函数
GkRpc.prototype.delete = function (recvJson) {
    let funcName = recvJson["funcName"];
    const eventId = recvJson["eventId"];

    if (funcName in this.functions) {
        delete this.functions[recvJson["funcName"]]

        this.sendJson({ "type": this.TypeDelete, "msg": { "success": true, "msg": "删除成功" }, "eventId": eventId })
    } else {
        this.sendJson({ "type": this.TypeDelete, "msg": { "success": false, "msg": "函数不存在" }, "eventId": eventId })
    }
}

// 执行指定函数
GkRpc.prototype.run = function (recvJson) {
    let _this = this;
    let funcName = recvJson["funcName"];
    const eventId = recvJson["eventId"];
    let hasSend = false;

    if (funcName in this.functions) {
        try {
            // 这里第一个参数就是：sendResult，第二个参数就是需要使用的参数
            this.functions[funcName](function (response) {
                if (response) {
                    hasSend = true;
                    _this.sendJson({ "type": _this.TypeRun, "msg": { "success": true, "msg": response }, "eventId": eventId })
                } else {
                    hasSend = true;
                    _this.sendJson({ "type": _this.TypeRun, "msg": { "success": true, "msg": funcName + " 执行完毕" }, "eventId": eventId })
                }
            }, recvJson)
        } catch (error) {
            hasSend = true;
            this.sendJson({ "type": this.TypeRun, "msg": { "success": false, "msg": funcName + " 执行出错：" + error.message }, "eventId": eventId })
        } finally {
            if (!hasSend) {
                this.sendJson({ "type": this.TypeRun, "msg": { "success": false, "msg": funcName + " 执行完毕" }, "eventId": eventId })
            }
        }
    } else {
        this.sendJson({ "type": this.TypeRun, "msg": { "success": false, "msg": "函数不存在" }, "eventId": eventId })
    }
}

// runEval 执行 eval
GkRpc.prototype.runEval = function (recvJson) {
    let js = recvJson["js"];
    const eventId = recvJson["eventId"];

    try {
        let result = eval.bind(window)(js);
        this.sendJson({ "type": this.TypeRun, "msg": { "success": true, "msg": result || "" }, "eventId": eventId })
    } catch (error) {
        this.sendJson({ "type": this.TypeRun, "msg": { "success": false, "msg": "执行出错：" + error.message }, "eventId": eventId })
    }
}

// send 发送信息
GkRpc.prototype.sendJson = function (msg) {
    console.log("发送数据：", msg);
    this.socket.send(JSON.stringify(msg))
}

window.demoRpc = new GkRpc("ws://127.0.0.1:8080/rpc")
