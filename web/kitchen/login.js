import { checkAuth } from "./js/auth.js";

async function login(data) {
  const resp = await fetch("/api/auth/login", {
    method: "POST",
    body: JSON.stringify(data),
  });
  if (!resp.ok) {
    throw Error("bad resp: ", resp.status);
  }

  const authed = await checkAuth();
  if (authed) {
    window.location.href = "/kitchen/";
  }
}

const form = document.querySelector("form#login-form");
form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const formRaw = new FormData(form);
  const loginData = Object.fromEntries(formRaw);
  await login(loginData);
});

const authed = await checkAuth();
if (authed) {
  window.location.href = "/kitchen/";
}
