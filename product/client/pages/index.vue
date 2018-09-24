<template>
  <section class="container">
    <div>
      <div>配信一覧</div>
      <div v-for="live in lives" :key="live.st">
        {{live}}
        <a :href="`/lives/${live.name}`">ライブを見る</a>
      </div>
    </div>
    <div>
      <input type="text" v-model="username">
      <input type="password" v-model="password">
      <button @click="login">ログイン</button>
    </div>
  </section>
</template>

<script>
export default {
  data () {
    return {
      username: '',
      password: ''
    }
  },
  async asyncData({ app }) {
    const url = process.server?'http://live-server:1323/api/lives':'/api/lives';
    const { data } = await app.$axios.get(url)
    return { lives: data }
  },
  methods:{
    login() {
      this.$axios.$post('http://localhost:8080/api/login', {username: this.username, password: this.password})
    }
  }
}
</script>

<style>
.container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.title {
  font-family: "Quicksand", "Source Sans Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; /* 1 */
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}
</style>

