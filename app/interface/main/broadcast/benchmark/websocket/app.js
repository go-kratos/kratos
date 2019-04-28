   
(function(win) {
  const rawHeaderLen = 18;
  const packetOffset = 0;
  const headerOffset = 4;
  const verOffset = 6;
  const opOffset = 8;
  const seqOffset = 12;
  const compressOffset = 16;
  const ctypeOffset = 17;

  var Client = function(options) {
      var MAX_CONNECT_TIMES = 10;
      var DELAY = 15000;
      this.options = options || {};
      this.createConnect(MAX_CONNECT_TIMES, DELAY);
  }

  Client.prototype.createConnect = function(max, delay) {
      var self = this;
      if (max === 0) {
          return;
      }
      connect();

      var textDecoder = new TextDecoder();
      var textEncoder = new TextEncoder();
      var heartbeatInterval;
      function connect() {
          var ws = new WebSocket('ws://172.22.33.126:7822/sub');
          ws.binaryType = 'arraybuffer';
          ws.onopen = function() {
            console.log("auth start")
              auth();
              register();
          }

          ws.onmessage = function(evt) {
              var data = evt.data;
              var dataView = new DataView(data, 0);
              var packetLen = dataView.getInt32(packetOffset);
              var headerLen = dataView.getInt16(headerOffset);
              var ver = dataView.getInt16(verOffset);
              var op = dataView.getInt32(opOffset);
              var seq = dataView.getInt32(seqOffset);
              var msgBody = textDecoder.decode(data.slice(headerLen, packetLen));

              console.log("receiveHeader: packetLen=" + packetLen, "headerLen=" + headerLen, "ver=" + ver, "op=" + op, "seq=" + seq,"body="+msgBody);

              switch(op) {
                  case 8:
                      // heartbeat
                      heartbeat();
                      heartbeatInterval = setInterval(heartbeat, 30 * 1000);
                  break;
                  case 3:
                      // heartbeat reply
                      console.log("receive: heartbeat online=", dataView.getInt32(rawHeaderLen));
                  break;
                  case 5:
                      // batch message
                      for (var offset=0; offset<data.byteLength; offset+=packetLen) {
                          // parse
                          var packetLen = dataView.getInt32(offset);
                          var headerLen = dataView.getInt16(offset+headerOffset);
                          var ver = dataView.getInt16(offset+verOffset);
                          var msgBody = textDecoder.decode(data.slice(offset+headerLen, offset+packetLen));
                          // callback
                          messageReceived(ver, msgBody);
                      }
                  break;
              }
          }

          ws.onclose = function() {
            console.log("closed")
              if (heartbeatInterval) clearInterval(heartbeatInterval);
              setTimeout(reConnect, delay);
          }

          function heartbeat() {
              var headerBuf = new ArrayBuffer(rawHeaderLen);
              var headerView = new DataView(headerBuf, 0);
              headerView.setInt32(packetOffset, rawHeaderLen);
              headerView.setInt16(headerOffset, rawHeaderLen);
              headerView.setInt16(verOffset, 1);
              headerView.setInt32(opOffset, 2);
              headerView.setInt32(seqOffset, 1);
              headerView.setInt8(compressOffset, 0);
              headerView.setInt8(ctypeOffset, 0);
              ws.send(headerBuf);
              console.log("send: heartbeat");
          }

          function auth() {
              var token ='{"room_id":"test://room_001","platform":"web"}' // ,"device_id":"123"
              var headerBuf = new ArrayBuffer(rawHeaderLen);
              var headerView = new DataView(headerBuf, 0);
              var bodyBuf = textEncoder.encode(token);
              headerView.setInt32(packetOffset, rawHeaderLen + bodyBuf.byteLength);
              headerView.setInt16(headerOffset, rawHeaderLen);
              headerView.setInt16(verOffset, 1);
              headerView.setInt32(opOffset, 7);
              headerView.setInt32(seqOffset, 1);
              headerView.setInt8(compressOffset, 0);
              headerView.setInt8(ctypeOffset, 0);
              ws.send(mergeArrayBuffer(headerBuf, bodyBuf));
          }

          function register() {
            var token ='{"operations":[1001,1002,1003]}' // ,"device_id":"123"
            var headerBuf = new ArrayBuffer(rawHeaderLen);
            var headerView = new DataView(headerBuf, 0);
            var bodyBuf = textEncoder.encode(token);
            headerView.setInt32(packetOffset, rawHeaderLen + bodyBuf.byteLength);
            headerView.setInt16(headerOffset, rawHeaderLen);
            headerView.setInt16(verOffset, 1);
            headerView.setInt32(opOffset, 14);
            headerView.setInt32(seqOffset, 3);
            headerView.setInt8(compressOffset, 0);
            headerView.setInt8(ctypeOffset, 0);
            ws.send(mergeArrayBuffer(headerBuf, bodyBuf));
        }

          function messageReceived(ver, body) {
              var notify = self.options.notify;
              if(notify) notify(body);
              console.log("messageReceived:", "ver=" + ver, "body=" + body);
          }

          function mergeArrayBuffer(ab1, ab2) {
              var u81 = new Uint8Array(ab1),
                  u82 = new Uint8Array(ab2),
                  res = new Uint8Array(ab1.byteLength + ab2.byteLength);
              res.set(u81, 0);
              res.set(u82, ab1.byteLength);
              return res.buffer;
          }

          function char2ab(str) {
              var buf = new ArrayBuffer(str.length);
              var bufView = new Uint8Array(buf);
              for (var i=0; i<str.length; i++) {
                  bufView[i] = str[i];
              }
              return buf;
          }

      }

      function reConnect() {
          self.createConnect(--max, delay * 2);
      }
  }

  win['MyClient'] = Client;
})(window);
