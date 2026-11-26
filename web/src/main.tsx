import { createRoot } from 'react-dom/client';
import { Entrada } from './app/entrada';
import './shared/styles/base.css';

const raiz = document.getElementById('raiz');
if (!raiz) {
    throw new Error('Elemento de montagem #raiz ausente no HTML.');
}

createRoot(raiz).render(<Entrada />);
