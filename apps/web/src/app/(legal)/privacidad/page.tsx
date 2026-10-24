export default function PrivacidadPage() {
  return (
    <article>
      <h1 className="text-display-3 mb-2">Política de Privacidad</h1>
      <p className="text-sm text-muted-foreground mb-12">Versión 1.1 | Vigente desde la fecha de publicación</p>

      <h2 className="text-heading-lg mt-10 mb-4">1. Nuestro doble rol</h2>
      <p>En UVICO tratamos datos personales bajo dos roles distintos:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li><strong>Como Responsable:</strong> Tratamos los datos de los Comercios (nombre, RIF, correo, facturación) y de los visitantes de uvic0.shop para prestar el servicio, facturar y mejorar la plataforma.</li>
        <li><strong>Como Encargado:</strong> Tratamos los datos de los Compradores (nombre, cédula, dirección, historial de pedidos) siguiendo las instrucciones del Comercio, quien decide sobre ellos. El Comercio es el responsable de esa data.</li>
      </ul>

      <h2 className="text-heading-lg mt-10 mb-4">2. Datos que recopilamos y su finalidad</h2>
      <p>Recopilamos únicamente los datos necesarios para operar de forma legal y segura:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li><strong>Registro y Facturación:</strong> Nombre, cédula o RIF, correo electrónico, teléfono y datos fiscales. Base legal: Ejecución de contrato y obligación legal (SENIAT).</li>
        <li><strong>Transaccionales:</strong> Datos de los pedidos para gestionar el dropshipping y la custodia de fondos. UVICO nunca almacena números de tarjetas de crédito; esos datos los maneja directamente la pasarela de pago.</li>
        <li><strong>Uso y Dispositivo:</strong> Dirección IP, navegador y páginas visitadas. Base legal: Interés legítimo para prevenir fraudes y mejorar la UX.</li>
        <li><strong>Verificación (KYC):</strong> Documentos de identidad. Base legal: Cumplimiento normativo (prevención de legitimación de capitales).</li>
      </ul>

      <h2 className="text-heading-lg mt-10 mb-4">3. Con quién compartimos los datos</h2>
      <p>No vendemos datos personales. Los compartimos estrictamente con:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li>Pasarelas de pago (para procesar las transacciones).</li>
        <li>Proveedores de dropshipping verificados (solo los datos necesarios para el despacho: nombre, dirección, teléfono).</li>
        <li>Proveedores de infraestructura (Supabase para base de datos y autenticación, Vercel para hosting).</li>
        <li>Autoridades competentes, cuando exista una orden legal válida.</li>
      </ul>

      <h2 className="text-heading-lg mt-10 mb-4">4. Conservación y Seguridad</h2>
      <p>Conservamos los datos mientras la cuenta esté activa y por los plazos legales exigidos (los soportes contables se conservan por 5 años según el Código de Comercio). Aplicamos cifrado en tránsito (TLS), control de accesos estricto y copias de seguridad.</p>

      <h2 className="text-heading-lg mt-10 mb-4">5. Tus derechos</h2>
      <p>Puedes solicitar acceso, rectificación, actualización o eliminación de tus datos escribiendo a <strong>privacidad@uvico.shop</strong>. Responderemos en un plazo máximo de 15 días hábiles. Si eres Comprador, dirige tu solicitud primero al Comercio; nosotros lo asistiremos.</p>
    </article>
  );
}
