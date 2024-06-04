<script setup>
import LoginForm from '../components/LoginForm.vue'
import HelloWorld from '../components/HelloWorld.vue'
</script>

<template>
  <div id="app">
  <header class="welcome">
    <img alt="Vue logo" class="logo" src="../components/icons/barry-wood.png" width="125" height="125" />

    <div class="wrapper">
      <HelloWorld msg="Welcome To social network!" />

      <nav>
        <RouterLink to="/">Login</RouterLink>
        <RouterLink to="/signup">SignUp</RouterLink>
      </nav>
    </div>
  </header><br>
    <LoginForm />

  </div>
</template>
<style scoped>
.welcome{
  width: 50rem;
}
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    border-right: 1px solid var(--color-border);
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
