<template>
  <section>
    <h1 v-if="!loggingIn">Login</h1>
    <div v-if="loggingIn" class="pacman">
      <img src="../assets/pacman_loading.svg" alt="" />
    </div>
    <div v-if="errorMessage" class="alert alert-danger" role="alert">
      {{ errorMessage }}
    </div>
    <form v-if="!loggingIn" @submit.prevent="login()">
      <div class="form-group col-md-6">
        <label for="email">Email</label>
        <input
          v-model="user.email"
          type="text"
          class="form-control"
          id="email"
          aria-describedby="emailHelp"
          placeholder="Enter a email"
          required
        />
        <h5 id="emailHelp" class="form-text text-muted">
          Enter a email to login.
        </h5>
      </div>
      <div class="form-group col-md-6">
        <label for="password">Password</label>
        <input
          v-model="user.password"
          type="password"
          class="form-control"
          id="password"
          aria-describedby="passwordHelp"
          placeholder="Enter a password"
          required
        />
        <h5 id="passwordHelp" class="form-text text-muted">
          Enter a password to login.
        </h5>
      </div>
      <div class="lbtns">
        <div class="login">
          <button type="submit" class="btn btn-primary btn-md">Login</button>
        </div>
        <div class="login">
          <router-link
            to="/signup"
            class="btn btn-primary btn-md"
            role="button"
          >
            Signup
          </router-link>
        </div>
      </div>
    </form>
  </section>
</template>

<script>
import Joi from "joi";
import axios from "axios";

const schema = Joi.object({
  email: Joi.string().email({ tlds: { allow: false } }),
  password: Joi.string().pattern(new RegExp("^[a-zA-Z0-9]{10,30}$")).trim(),
});

export default {
  data: () => ({
    errorMessage: "",
    loggingIn: false,
    user: {
      email: "",
      password: "",
    },
  }),
  methods: {
    login() {
      this.errorMessage = "";
      if (this.validUser()) {
        this.loggingIn = true;
        axios
          .post("/v1/login", {
            user: {
              email: this.user.email,
              password: this.user.password,
            },
          })
          .then((result) => {
            localStorage.token = result.data.token.plaintext;
            setTimeout(() => {
              this.loggingIn = false;
              this.$router.push("/calendar");
            }, 300);
          })
          .catch((err) => {
            setTimeout(() => {
              this.loggingIn = false;
              this.errorMessage = err.message;
            }, 300);
          });
      }
    },
    validUser() {
      const result = schema.validate(this.user);
      if (!Object.prototype.hasOwnProperty.call(result, "error")) {
        return true;
      }
      if (result.error.message.includes("email")) {
        this.errorMessage = "Email is invalid.";
      } else {
        this.errorMessage = "Password is invalid.";
      }

      return false;
    },
  },
};
</script>

<style>
.pacman {
  height: 550px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.lbtns {
  display: flex;
  justify-content: start;
  padding-left: 12px;
}
.login {
  margin-right: 25px;
}
</style>
