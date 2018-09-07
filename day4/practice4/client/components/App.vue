<template>
<div>
  <h1>どうがみようぜ</h1>
  <div>
    <div v-for="live in lives">
      <div class="live-name">{{live.Name}}</div>
      <button @click="startLive(live)">liveを見る</button>
    </div>
  </div>
  <video id="video"></video>
</div>
</template>

<script>
import Hls from 'hls.js'
import axios from 'axios'
axios.defaults.baseURL = 'http://localhost:8888/api'
export default {
  name: 'App',
  data() {
    return {
      lives: []
    }
  },
  async mounted () {
    const po = await axios.get('lives')
    console.log(po)
    this.lives = po.data
  },
  methods: {
    startLive(live) {
      const video = document.getElementById('video');
      if(Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(`http://localhost:8888/${live.Name}.m3u8`);
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED,function() {
          video.play();
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = `http://localhost:8888/${live.Name}.m3u8`;
        video.addEventListener('loadedmetadata',function() {
          video.play();
        });
      }
    }
  }
}
</script>

<style>

</style>

