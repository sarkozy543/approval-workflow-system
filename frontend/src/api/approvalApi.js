const API_BASE = "http://localhost:8080";

async function handleResponse(res) {
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed with status ${res.status}`);
  }
  if (res.status === 204) return null;
  return res.json();
}

export async function fetchRequests() {
  const res = await fetch(`${API_BASE}/requests`);
  return handleResponse(res);
}

export async function createRequest(body) {
  const res = await fetch(`${API_BASE}/requests`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User": "yusuf",
    },
    body: JSON.stringify(body),
  });
  return handleResponse(res);
}

export async function getRequest(id) {
  const res = await fetch(`${API_BASE}/requests/${id}`);
  return handleResponse(res);
}

export async function approveRequest(id, note) {
  const res = await fetch(`${API_BASE}/requests/${id}/approve`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User": "yusuf",
    },
    body: JSON.stringify({ note }),
  });
  return handleResponse(res);
}

export async function rejectRequest(id, note) {
  const res = await fetch(`${API_BASE}/requests/${id}/reject`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User": "yusuf",
    },
    body: JSON.stringify({ note }),
  });
  return handleResponse(res);
}

export async function fetchLogs(id) {
  const res = await fetch(`${API_BASE}/requests/${id}/logs`);
  return handleResponse(res);
}

