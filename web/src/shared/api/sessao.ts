export type Conta = {
    id: number;
    nome: string;
    email: string;
    perfil: string;
    ultimoAcesso: string | null;
};

let tokenAcesso = '';

async function renovarSessao(): Promise<boolean> {
    const resposta = await fetch('/api/auth/refresh', {
        method: 'POST',
        credentials: 'same-origin',
    });
    tokenAcesso = resposta.headers.get('Authorization')?.replace(/^Bearer\s+/i, '') ?? '';
    return resposta.ok && Boolean(tokenAcesso);
}

export async function consultarSessao(): Promise<Conta | null> {
    if (!await renovarSessao()) return null;
    return solicitar<Conta>('/api/auth/me');
}

export async function entrar(email: string, senha: string): Promise<Conta> {
    const resposta = await fetch('/api/auth/login', {
        method: 'POST',
        credentials: 'same-origin',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, senha }),
    });
    tokenAcesso = resposta.headers.get('Authorization')?.replace(/^Bearer\s+/i, '') ?? '';
    return lerResposta<Conta>(resposta);
}

export async function sair(): Promise<void> {
    try {
        await solicitar('/api/auth/logout', { method: 'POST' });
    } finally {
        tokenAcesso = '';
    }
}

async function solicitar<T>(caminho: string, opcoes: RequestInit = {}): Promise<T> {
    const fazerRequisicao = () => fetch(caminho, {
        ...opcoes,
        credentials: 'same-origin',
        headers: { ...opcoes.headers, Authorization: `Bearer ${tokenAcesso}` },
    });
    let resposta = await fazerRequisicao();
    if (resposta.status === 401 && await renovarSessao()) resposta = await fazerRequisicao();
    return lerResposta<T>(resposta);
}

async function lerResposta<T>(resposta: Response): Promise<T> {
    const corpo: unknown = await resposta.json().catch(() => null);
    if (!resposta.ok) {
        const mensagem = typeof corpo === 'object' && corpo !== null && 'erro' in corpo
            ? String(corpo.erro)
            : 'Não foi possível concluir a solicitação.';
        throw new Error(mensagem);
    }
    return corpo as T;
}
