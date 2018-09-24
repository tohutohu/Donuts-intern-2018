<template>
  <div>
    <div style="display:flex;">
      <div class="video-container">
        <video id="video" controls="on"></video>
      </div>
      <div class="chat-container">
        <div class="messages-container">
          <div v-for="m in messages" :key="m.text" class="message">
            <div class="username">{{m.username}}</div>
            <div class="text">{{m.text}}</div>
          </div>
        </div>
        <div class="input-component">
          <input type="text" @keydown.enter="sendMessage" v-model="messageText">
          <button @click="sendMessage"><i class="far fa-paper-plane"></i></button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Hls from 'hls.js'
import WS from '@/utils/WS.js'

export default {
  name: 'Live',
  data () {
    return {
      socket: null,
      messageText: '',
      messages: []
    }
  },
  async mounted () {
    this.socket = new WS()
    this.socket.on('message-create', data => {
      this.messages.push(data)
    })

    this.socket.on('initial-data', data => {
      this.messages = data.reverse()
    })
    await this.socket.connect('ws://localhost:8080/ws/' + this.$route.params.id)
    this.startVideo()
  },
  methods: {
    startVideo() {
      const video = document.getElementById('video');
      if(Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(`/hls/${this.$route.params.id}.m3u8`);
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED,function() {
          video.play();
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = `/hls/${this.$route.params.id}.m3u8`;
        video.addEventListener('loadedmetadata',function() {
          video.play();
        });
      }
    },
    sendMessage() {
      if (this.messageText === '') {
        return
      }

      this.socket.send('message-create', {text: this.messageText})
      this.messageText = ''
    }
  }
}
</script>


<style>
#video {
  width: 100%;
}

.video-container {
  width: 50%;
}

.chat-container {
  width: 480px;
  height: 600px;
  padding: 12px;
}

.messages-container {
  height: calc(100% - 60px);
  overflow: scroll;
}

.input-component {
  width: 100%;
  bottom: 0;
}

.input-component > input {
  width: calc(100% - 36px);
}

.input-component > button {
  width: 24px;
  height: 24px;
  border: none;
  padding: 2px;
  background: white;
}

.input-component > button > i{
  font-size: 20px;
}

.input-component > button:hover {
  background: #dadada;
  transition-duration: 0.2s;
}

.username {
  color: #dadada;
  margin-right: 8px;
}

.message {
  margin-bottom: 12px;
  display: flex;
  word-wrap: wrap;
  word-break: breakall;

}
</style>
