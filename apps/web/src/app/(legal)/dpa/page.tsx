export default function DPAPage() {
  return (
    <article>
      <h1 className="text-display-3 mb-2">Acuerdo de Tratamiento de Datos</h1>
      <p className="text-sm text-muted-foreground mb-12">Versión 1.1 | Vigente desde la fecha de publicación</p>

      <h2 className="text-heading-lg mt-10 mb-4">1. Partes y objeto</h2>
      <p>Este Acuerdo de Tratamiento de Datos (el “DPA”) forma parte integrante de los Términos y Condiciones de UVICO, S.A. (en adelante, “UVICO” o “el Encargado”) y regula el tratamiento de los datos personales que el Comercio (en adelante, “el Responsable”) le confía a UVICO para la operación de su tienda en línea.</p>
      <p>El Comercio actúa como responsable de los datos de sus Compradores. UVICO actúa exclusivamente como encargado del tratamiento, siguiendo las instrucciones del Comercio y utilizando los datos únicamente para prestar el servicio de plataforma de e-commerce y gestión de dropshipping.</p>

      <h2 className="text-heading-lg mt-10 mb-4">2. Naturaleza y finalidad del tratamiento</h2>
      <p>UVICO tratará los siguientes datos personales de los Compradores por cuenta del Comercio:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li><strong>Datos de identificación y contacto:</strong> Nombre, cédula o RIF, correo electrónico, número de teléfono y dirección de entrega.</li>
        <li><strong>Datos transaccionales:</strong> Historial de pedidos, productos adquiridos, montos y estado de los envíos.</li>
        <li><strong>Datos técnicos:</strong> Dirección IP y datos de navegación dentro de la tienda (necesarios para la prevención de fraudes y analíticas).</li>
      </ul>
      <p>La finalidad del tratamiento es estrictamente: (a) gestionar los pedidos y la logística de entrega; (b) procesar los pagos a través de las pasarelas; (c) gestionar la retención y liberación de fondos en el flujo de dropshipping; y (d) prevenir fraudes y legitimación de capitales.</p>

      <h2 className="text-heading-lg mt-10 mb-4">3. Obligaciones del Comercio (Responsable)</h2>
      <p>El Comercio garantiza que tiene una base legal legítima para recopilar y tratar los datos de sus Compradores, y que ha informado a estos sobre su propia Política de Privacidad. El Comercio se compromete a no instruir a UVICO a realizar tratamientos ilícitos o que violen la privacidad de los usuarios.</p>

      <h2 className="text-heading-lg mt-10 mb-4">4. Obligaciones de UVICO (Encargado)</h2>
      <p>UVICO se obliga a:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li>Tratar los datos personales únicamente para las finalidades descritas y bajo las instrucciones del Comercio.</li>
        <li>Mantener la confidencialidad absoluta. Todo el personal de UVICO con acceso a los datos está sujeto a acuerdos de confidencialidad.</li>
        <li>No vender, alquilar ni compartir los datos con terceros no autorizados.</li>
        <li>Asistir al Comercio en la respuesta a las solicitudes de los Compradores (acceso, rectificación o eliminación de sus datos).</li>
        <li>Notificar al Comercio sin demora indebida (y en un plazo máximo de 72 horas) cualquier brecha de seguridad que afecte los datos personales.</li>
      </ul>

      <h2 className="text-heading-lg mt-10 mb-4">5. Subencargados del tratamiento</h2>
      <p>El Comercio autoriza a UVICO a contratar subencargados para la infraestructura del servicio. UVICO utiliza proveedores de nivel empresarial que cumplen con altos estándares de seguridad:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li><strong>Supabase Inc.:</strong> Para el almacenamiento seguro de la base de datos y la autenticación de usuarios.</li>
        <li><strong>Vercel Inc.:</strong> Para el alojamiento (hosting) y despliegue del frontend.</li>
        <li><strong>Pasarelas de Pago (Wayu Pay, VeriPay, etc.):</strong> Actúan como responsables independientes del tratamiento de los datos financieros (tarjetas, cuentas bancarias) y no son subencargados de UVICO.</li>
      </ul>

      <h2 className="text-heading-lg mt-10 mb-4">6. Medidas de seguridad</h2>
      <p>UVICO implementa medidas técnicas y organizativas de nivel bancario, incluyendo: cifrado de datos en tránsito (TLS/SSL), control de acceso basado en roles (RBAC), políticas de contraseñas estrictas, copias de seguridad automatizadas y monitoreo continuo de la infraestructura. UVICO no almacena números completos de tarjetas de crédito.</p>

      <h2 className="text-heading-lg mt-10 mb-4">7. Fin del servicio y eliminación de datos</h2>
      <p>Al terminar la relación contractual (por cancelación o suspensión), UVICO otorgará al Comercio un plazo de 30 días calendario para exportar sus datos. Vencido este plazo, UVICO procederá a eliminar o anonimizar los datos personales de los Compradores, salvo aquellos que deba conservar por obligaciones legales, contables o tributarias (como los registros de facturación exigidos por el SENIAT).</p>
    </article>
  );
}
