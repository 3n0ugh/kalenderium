<template>
  <section>
    <div class="home text-center">
      <div class="jumbotron">
        <img
          src="../assets/kalenderium_logo.png"
          alt=""
          width="300"
          height="300"
        />
        <p class="lead mid">Javascript is giving me pain but, it's also fun.</p>
        <hr class="my-4" />
        <p>Another hard-to-understand project.</p>
        <div class="lb-tns">
          <div class="login">
            <router-link to="/signup" class="btn btn-primary" role="button">
              Signup
            </router-link>
          </div>
          <div class="login">
            <router-link
              to="/login"
              class="btn btn-primary btn-margin"
              role="button"
            >
              Login
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style>
.mid {
  margin-top: 40px;
}
.lb-tns {
  display: flex;
  justify-content: start;
  padding-left: 42%;
}
.login {
  margin-right: 25px;
}
</style>
