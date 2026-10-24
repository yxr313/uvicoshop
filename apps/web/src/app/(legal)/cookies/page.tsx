export default function CookiesPage() {
  return (
    <article>
      <h1 className="text-display-3 mb-2">Política de Cookies</h1>
      <p className="text-sm text-muted-foreground mb-12">Versión 1.1 | Vigente desde la fecha de publicación</p>

      <h2 className="text-heading-lg mt-10 mb-4">1. Qué son las cookies y por qué las usamos</h2>
      <p>Las cookies son pequeños archivos de texto que un sitio web guarda en tu dispositivo. Las usamos para mantener tu sesión iniciada, recordar tus preferencias y entender cómo interactúas con la plataforma para mejorar el producto.</p>

      <h2 className="text-heading-lg mt-10 mb-4">2. Categorías de cookies</h2>
      <p>Utilizamos las siguientes categorías. Las duraciones son máximas y se miden desde la última visita.</p>
      
      <div className="overflow-x-auto my-8">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-border">
              <th className="py-3 px-4 font-semibold text-foreground">Categoría</th>
              <th className="py-3 px-4 font-semibold text-foreground">Finalidad</th>
              <th className="py-3 px-4 font-semibold text-foreground">Duración</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            <tr>
              <td className="py-3 px-4 text-sm font-medium">Estrictamente necesarias</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Iniciar sesión, mantener la sesión, seguridad del carrito y funcionamiento general.</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Sesión a 12 meses</td>
            </tr>
            <tr>
              <td className="py-3 px-4 text-sm font-medium">Preferencias</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Recordar idioma, moneda local (VES/USD) y configuración de la interfaz.</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Hasta 12 meses</td>
            </tr>
            <tr>
              <td className="py-3 px-4 text-sm font-medium">Analíticas</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Medir el uso del sitio (páginas visitadas, tiempo) para mejorar la UX.</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Hasta 24 meses</td>
            </tr>
            <tr>
              <td className="py-3 px-4 text-sm font-medium">Marketing</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Medir campañas publicitarias y mostrar anuncios relevantes.</td>
              <td className="py-3 px-4 text-sm text-muted-foreground">Hasta 90 días</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h2 className="text-heading-lg mt-10 mb-4">3. Consentimiento y Control</h2>
      <p>Al ingresar por primera vez te mostramos un aviso para aceptar todas, rechazar las no esenciales o configurarlas. Las cookies estrictamente necesarias no requieren consentimiento. Puedes cambiar tu decisión en cualquier momento desde el enlace "Configuración de cookies" en el pie del sitio o bloqueándolas directamente en tu navegador (Chrome, Safari, Firefox).</p>
    </article>
  );
}
