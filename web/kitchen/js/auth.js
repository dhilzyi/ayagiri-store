export async function checkAuth() {
  const res = await fetch("/api/auth/me", {
    credentials: "include",
  });
  return res.ok;
}
