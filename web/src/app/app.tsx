import marca from './mira.txt?raw';
import './app.css';

function App() {
    return (
        <main className="tela-inicial">
            <div className="painel-inicial">
                <pre className="marca-inicial" aria-hidden="true">{marca.trimEnd()}</pre>
                <h1>
                    Mira<span aria-hidden="true">.</span>
                </h1>
                <p>A autenticação do console está em desenvolvimento.</p>
            </div>
        </main>
    );
}

export { App };
