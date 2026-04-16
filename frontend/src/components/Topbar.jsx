export default function Topbar() {
  return (
    <header className="topbar">
      <input
        type="search"
        className="search"
        placeholder="Procurar criadoras, topicos e posts"
      />
      <div className="top-actions">
        <button className="ghost-btn">Notificacoes</button>
        <button className="ghost-btn">Config</button>
      </div>
    </header>
  );
}
