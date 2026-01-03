import { useEffect, useState } from "react";
import {
  fetchRequests,
  createRequest,
  approveRequest,
  rejectRequest,
  fetchLogs,
} from "./api/approvalApi";
import "./App.css";

function App() {
  const [requests, setRequests] = useState([]);
  const [selected, setSelected] = useState(null);
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState(null);

  const [form, setForm] = useState({
    source_env: "dev",
    target_env: "staging",
    requested_by: "yusuf",
    change_payload: {
      flag: "new_feature",
      enabled: true,
    },
  });

  useEffect(() => {
    loadRequests();
  }, []);

  async function loadRequests() {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchRequests();
      setRequests(data);
      if (!selected && data.length > 0) {
        handleSelect(data[0]);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleSelect(req) {
    setSelected(req);
    try {
      const data = await fetchLogs(req.id);
      setLogs(data);
    } catch (err) {
      console.error("Failed to load logs", err);
      setLogs([]);
    }
  }

  async function handleCreate(e) {
    e.preventDefault();
    try {
      setCreating(true);
      setError(null);

      const body = {
        source_env: form.source_env,
        target_env: form.target_env,
        requested_by: form.requested_by,
        change_payload: form.change_payload,
      };

      const created = await createRequest(body);
      await loadRequests();
      setSelected(created);
      const data = await fetchLogs(created.id);
      setLogs(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setCreating(false);
    }
  }

  async function handleApprove() {
    if (!selected) return;
    const note = window.prompt("Onay notu (opsiyonel):", "") || "";

    try {
      setLoading(true);
      setError(null);
      const updated = await approveRequest(selected.id, note);
      setSelected(updated);
      await loadRequests();
      const data = await fetchLogs(updated.id);
      setLogs(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleReject() {
    if (!selected) return;
    const note = window.prompt("Red notu (opsiyonel):", "") || "";

    try {
      setLoading(true);
      setError(null);
      const updated = await rejectRequest(selected.id, note);
      setSelected(updated);
      await loadRequests();
      const data = await fetchLogs(updated.id);
      setLogs(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Approval Workflow Panel</h1>
        <button className="reload-btn" onClick={loadRequests} disabled={loading}>
          Yenile
        </button>
      </header>

      {error && <div className="error-banner">⚠ {error}</div>}

      <div className="app-body">
        {/* SOL: LİSTE + YENİ FORM */}
        <div className="sidebar">
          <section className="card">
            <h2>Requests</h2>
            {loading && <p>Yükleniyor...</p>}
            {!loading && requests.length === 0 && <p>Hiç request yok.</p>}
            <ul className="request-list">
              {requests.map((r) => (
                <li
                  key={r.id}
                  className={
                    selected && selected.id === r.id
                      ? "request-item selected"
                      : "request-item"
                  }
                  onClick={() => handleSelect(r)}
                >
                  <div className="req-envs">
                    <span>{r.source_env}</span>
                    <span className="arrow">→</span>
                    <span>{r.target_env}</span>
                  </div>
                  <div className={`status badge ${r.status.toLowerCase()}`}>
                    {r.status}
                  </div>
                  <div className="req-meta">
                    <small>{r.requested_by}</small>
                  </div>
                </li>
              ))}
            </ul>
          </section>

          <section className="card">
            <h2>Yeni Request</h2>
            <form onSubmit={handleCreate} className="form">
              <label>
                Source Env
                <input
                  value={form.source_env}
                  onChange={(e) =>
                    setForm({ ...form, source_env: e.target.value })
                  }
                />
              </label>
              <label>
                Target Env
                <input
                  value={form.target_env}
                  onChange={(e) =>
                    setForm({ ...form, target_env: e.target.value })
                  }
                />
              </label>
              <label>
                Requested By
                <input
                  value={form.requested_by}
                  onChange={(e) =>
                    setForm({ ...form, requested_by: e.target.value })
                  }
                />
              </label>
              <label>
                Flag Name
                <input
                  value={form.change_payload.flag}
                  onChange={(e) =>
                    setForm({
                      ...form,
                      change_payload: {
                        ...form.change_payload,
                        flag: e.target.value,
                      },
                    })
                  }
                />
              </label>
              <label className="checkbox">
                <input
                  type="checkbox"
                  checked={form.change_payload.enabled}
                  onChange={(e) =>
                    setForm({
                      ...form,
                      change_payload: {
                        ...form.change_payload,
                        enabled: e.target.checked,
                      },
                    })
                  }
                />
                Enabled
              </label>
              <button type="submit" disabled={creating}>
                {creating ? "Oluşturuluyor..." : "Oluştur"}
              </button>
            </form>
          </section>
        </div>

        {/* SAĞ: DETAY + LOGS */}
        <div className="main">
          <section className="card">
            <h2>Detay</h2>
            {!selected && <p>Soldan bir request seç.</p>}
            {selected && (
              <div>
                <p>
                  <strong>ID:</strong> {selected.id}
                </p>
                <p>
                  <strong>Env:</strong> {selected.source_env} →{" "}
                  {selected.target_env}
                </p>
                <p>
                  <strong>Status:</strong>{" "}
                  <span className={`badge ${selected.status.toLowerCase()}`}>
                    {selected.status}
                  </span>
                </p>
                <p>
                  <strong>Requested By:</strong> {selected.requested_by}
                </p>
                <p>
                  <strong>Change Payload:</strong>
                  <pre className="payload">
                    {JSON.stringify(selected.change_payload, null, 2)}
                  </pre>
                </p>

                <div className="actions">
                  <button
                    onClick={handleApprove}
                    disabled={loading || selected.status !== "PENDING"}
                  >
                    Approve
                  </button>
                  <button
                    onClick={handleReject}
                    disabled={loading || selected.status !== "PENDING"}
                  >
                    Reject
                  </button>
                </div>
              </div>
            )}
          </section>

          <section className="card">
            <h2>Logs</h2>
            {logs.length === 0 && <p>Bu request için log yok.</p>}
            {logs.length > 0 && (
              <ul className="logs">
                {logs.map((log) => (
                  <li key={log.id} className="log-item">
                    <div className="log-main">
                      <span className={`badge small ${log.action.toLowerCase()}`}>
                        {log.action}
                      </span>
                      <span className="log-by">{log.action_by}</span>
                      <span className="log-date">
                        {new Date(log.created_at).toLocaleString()}
                      </span>
                    </div>
                    {log.note && <div className="log-note">{log.note}</div>}
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}

export default App;
