import Vue from "vue";
import VueRouter from "vue-router";
import HomeView from "../views/HomeView.vue";
import CalendarView from "@/views/CalendarView";
import SignupView from "@/views/SignupView";
import LoginView from "@/views/LoginView";

Vue.use(VueRouter);

function loggedInRedirectDashboard(to, from, next) {
  if (localStorage.token) {
    next("/calendar");
  } else {
    next();
  }
}

function isLoggedIn(to, from, next) {
  if (localStorage.token) {
    next();
  } else {
    next("/login");
  }
}

const routes = [
  {
    path: "/",
    name: "home",
    component: HomeView,
  },
  {
    path: "/calendar",
    name: "calendar",
    component: CalendarView,
    beforeEnter: isLoggedIn,
  },
  {
    path: "/signup",
    name: "signup",
    component: SignupView,
    beforeEnter: loggedInRedirectDashboard,
  },
  {
    path: "/login",
    name: "login",
    component: LoginView,
    beforeEnter: loggedInRedirectDashboard,
  },
];

const router = new VueRouter({
  routes,
});

export default router;
