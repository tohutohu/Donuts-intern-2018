export default class WS{
  constructor() {
    this.cb = {}
  }

  connect(url) {
    return new Promise((res) => {
      this.socket = new WebSocket(url)
      this.url = url
      this.socket.onopen = data => {
        console.log('websocket connected to ' + url)
        res()
      }

      this.socket.onmessage = ({ data }) => {
        const payload = JSON.parse(data)
        console.log('websocket received from ' + this.url, payload)
        if (Array.isArray(this.cb[payload.type])) {
          this.cb[payload.type].forEach(f => {
            f(payload.data)
          })
        }
      }
      this.socket.onclose = () => {
        console.log('websocket connection closed')
        setTimeout(() => {
          console.log('try reconnection...')
          this.connect(this.url)
        }, 1000)
      }
    })
  }

  send(type, data) {
    this.socket.send(JSON.stringify({type: type, data:data}))
  }

  on(type, f) {
    if (!Array.isArray(this.cb[type])){
      this.cb[type] = []
    }

    this.cb[type].push(f)
  }
}
