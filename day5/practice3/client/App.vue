<template>
  <div>
    <div>test po</div>
    <div>
      <div v-for="m in messages" :key="m.data.text">
        {{m.data.text}}
      </div>
    </div>

    <input type="text" v-model="text">
    <button @click="sendMessage">送信</button>
  
  </div>
</template>


<script>
export default {
  name: 'App',
  data () {
    return {
      ws: null,
      text: '',
      messages: []
    }
  },
  mounted() {
	  console.log('mounted')
    if(location.hash != '') {
      this.ws = new WebSocket('ws://127.0.0.1:8888/ws/test')
    } else {
      this.ws = new WebSocket('ws://127.0.0.1:1323/ws/test')
    }
    this.ws.onmessage = data => {
      const message = JSON.parse(data.data)
      console.log(message)
      if (message.type === 'message') {
        this.messages.push(message)
      }
    }

    this.ws.onopen = data => {
      console.log(data)
    }

    this.ws.onclose = data => {
      console.log(data)
    }
  },
  methods: {
    sendMessage() {
      this.ws.send(JSON.stringify({type:'message', data: {text: this.text}}))
    }
  }
}
</script>
