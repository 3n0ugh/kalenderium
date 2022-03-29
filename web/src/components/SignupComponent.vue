<template>
  <section>
    <h1>Signup</h1>
    <div v-if="signingUp" class="pacman">
      <img src="../assets/pacman_loading.svg" />
    </div>
    <div v-if="errorMessage" class="alert alert-danger" role="alert">
      {{ errorMessage }}
    </div>
    <form v-if="!signingUp" @submit.prevent="signup">
      <div class="form-row">
        <div class="form-group col-md-6">
          <label for="email">Email</label>
          <input
            v-model="user.email"
            type="text"
            class="form-control"
            id="email"
            placeholder="Enter a email"
            required
          />
        </div>
        <div class="form-group col-md-6">
          <label for="password">Password</label>
          <input
            v-model="user.password"
            type="password"
            class="form-control"
            id="password"
            aria-describedby="passwordHelp"
            placeholder="Password"
            required
          />
          <h5 id="passwordHelp" class="form-text text-muted">
            Password must be longer than 10.
          </h5>
        </div>

        <div class="form-group col-md-6">
          <label for="confirmPassword">Confirm Password</label>
          <input
            v-model="user.confirmPassword"
            type="password"
            class="form-control"
            id="confirmPassword"
            aria-describedby="confirmPasswordHelp"
            placeholder="Password"
            required
          />
          <h5 id="confirmPasswordHelp" class="form-text text-muted">
            Please confirm your password.
          </h5>
        </div>
      </div>
      <button type="submit" class="btn btn-primary col-md-1 l-btn">
        Signup
      </button>
    </form>
  </section>
</template>

<script>
import Joi from "joi";
import axios from "axios";

const schema = Joi.object({
  email: Joi.string().email({ tlds: { allow: false } }),
  password: Joi.string().pattern(new RegExp("^[a-zA-Z0-9]{10,30}$")).trim(),
  confirmPassword: Joi.string()
    .pattern(new RegExp("^[a-zA-Z0-9]{10,30}$"))
    .trim(),
});

export default {
  data: () => ({
    signingUp: false,
    errorMessage: "",
    user: {
      email: "",
      password: "",
      confirmPassword: "",
    },
  }),
  watch: {
    user: {
      handler() {
        this.errorMessage = "";
      },
      deep: true,
    },
  },
  methods: {
    signup() {
      this.errorMessage = "";
      if (this.validUser()) {
        this.signingUp = true;
        axios
          .post("/v1/signup", {
            user: {
              email: this.user.email,
              password: this.user.password,
            },
          })
          .then((result) => {
            localStorage.token = result.data.token;
            setTimeout(() => {
              this.signingUp = false;
              this.$router.push("/calendar");
            }, 700);
          })
          .catch((err) => {
            setTimeout(() => {
              this.signingUp = false;
              this.errorMessage = err.message;
            }, 700);
          });
      }
    },
    validUser() {
      if (this.user.password !== this.user.confirmPassword) {
        this.errorMessage = "Passwords must match.";
        return false;
      }
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
.l-btn {
  margin-left: 12px;
}
</style>
