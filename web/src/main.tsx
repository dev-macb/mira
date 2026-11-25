import { createRoot } from 'react-dom/client';
import { App } from './app/app';
import './shared/styles/base.css';

const raiz = document.getElementById('raiz');
if (!raiz) {
    throw new Error('Elemento de montagem #raiz ausente no HTML.');
}

createRoot(raiz).render(<App />);
