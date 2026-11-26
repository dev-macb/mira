import { useEffect, useRef, useState, type FormEvent } from 'react';
import { consultarSessao, entrar, sair, type Conta } from '../shared/api/sessao';
import marca from './mira.txt?raw';
import './entrada.css';

function Entrada() {
    const [conta, definirConta] = useState<Conta | null>(null);
    const [carregando, definirCarregando] = useState(true);
    const [enviando, definirEnviando] = useState(false);
    const [mostrarSenha, definirMostrarSenha] = useState(false);
    const [erro, definirErro] = useState('');
    const campoEmail = useRef<HTMLInputElement>(null);
    const titulo = useRef<HTMLHeadingElement>(null);

    useEffect(() => {
        consultarSessao()
            .then(definirConta)
            .catch(() => definirConta(null))
            .finally(() => definirCarregando(false));
    }, []);

    useEffect(() => {
        if (carregando) return;
        if (conta) titulo.current?.focus();
        else campoEmail.current?.focus();
    }, [carregando, conta]);

    async function enviarFormulario(evento: FormEvent<HTMLFormElement>) {
        evento.preventDefault();
        if (enviando) return;
        definirEnviando(true);
        definirErro('');
        const dados = new FormData(evento.currentTarget);
        try {
            definirConta(await entrar(String(dados.get('email')), String(dados.get('senha'))));
        } catch (falha) {
            definirErro(falha instanceof Error ? falha.message : 'Não foi possível entrar.');
            campoEmail.current?.focus();
        } finally {
            definirEnviando(false);
        }
    }

    async function encerrarSessao() {
        try { await sair(); } finally { definirConta(null); }
    }

    if (carregando) return <main className="tela-acesso" aria-busy="true"><p>Carregando…</p></main>;

    if (conta) return (
        <main className="tela-acesso">
            <section className="painel-acesso painel-conta" aria-labelledby="titulo-conta">
                <pre className="marca-acesso" aria-hidden="true">{marca.trimEnd()}</pre>
                <h1 id="titulo-conta" ref={titulo} tabIndex={-1}>Mira</h1>
                <p>Olá, {conta.nome}.</p>
                <p className="texto-secundario">A autenticação está pronta. Os projetos entram na próxima etapa.</p>
                <button className="botao-primario" type="button" onClick={encerrarSessao}>Sair da conta</button>
            </section>
        </main>
    );

    return (
        <main className="tela-acesso">
            <section className="painel-acesso" aria-labelledby="titulo-acesso">
                <pre className="marca-acesso" aria-hidden="true">{marca.trimEnd()}</pre>
                <h1 id="titulo-acesso">Mira</h1>
                <p className="texto-secundario">Acesso ao console</p>
                <form className="formulario-acesso" onSubmit={enviarFormulario}>
                    <label htmlFor="email">E-mail</label>
                    <input id="email" ref={campoEmail} name="email" type="email" autoComplete="username" required maxLength={254} onChange={() => definirErro('')} />
                    <label htmlFor="senha">Senha</label>
                    <div className="campo-senha">
                        <input id="senha" name="senha" type={mostrarSenha ? 'text' : 'password'} autoComplete="current-password" required onChange={() => definirErro('')} />
                        <button className="botao-senha" type="button" aria-label={mostrarSenha ? 'Ocultar senha' : 'Mostrar senha'} onClick={() => definirMostrarSenha(!mostrarSenha)}>{mostrarSenha ? 'Ocultar' : 'Mostrar'}</button>
                    </div>
                    {erro && <p className="erro-acesso" role="alert">{erro}</p>}
                    <button className="botao-primario" type="submit" disabled={enviando}>{enviando ? 'Entrando…' : 'Entrar'}</button>
                </form>
            </section>
        </main>
    );
}

export { Entrada };
