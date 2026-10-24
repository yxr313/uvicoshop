export default function TerminosPage() {
  return (
    <article>
      <h1 className="text-display-3 mb-2">Términos y Condiciones</h1>
      <p className="text-sm text-muted-foreground mb-12">Versión 1.1 | Vigente desde la fecha de publicación</p>

      <h2 className="text-heading-lg mt-10 mb-4">1. Quiénes somos y aceptación</h2>
      <p>Estos Términos y Condiciones (los “Términos”) regulan el acceso y uso de la plataforma UVICO, disponible en www.uvico.shop y en sus aplicaciones asociadas (la “Plataforma”), operada por UVICO, S.A. (en adelante, “UVICO”), sociedad mercantil venezolana.</p>
      <p>Al crear una cuenta, marcar la casilla de aceptación o usar la Plataforma, el Comercio declara que leyó, entendió y acepta estos Términos, junto con la Política de Privacidad, la Política de Cookies, la Política de Uso Aceptable, la Política de Reembolsos y Cancelación y el Acuerdo de Tratamiento de Datos (el “DPA”).</p>

      <h2 className="text-heading-lg mt-10 mb-4">2. Descripción del servicio y el modelo de custodia</h2>
      <p>UVICO es una plataforma de software como servicio (SaaS) y un Ecommerce Builder. El servicio permite a emprendedores y empresas (los “Comercios”) crear y administrar tiendas en línea.</p>
      <p><strong>El modelo de Dropshipping Verificado:</strong> UVICO integra una red de proveedores verificados. Cuando un Comercio vende un producto de esta red, UVICO actúa como custodio de los fondos. El dinero del Comprador es procesado por la Pasarela de Pago, retenido por UVICO para protección de fraudes, y utilizado para pagar directamente al proveedor. UVICO no custodia los fondos de productos de inventario propio del Comercio.</p>
      <p><strong>Comisión:</strong> Por cada venta de dropshipping, UVICO cobrará una comisión del cinco por ciento (5%) sobre el monto total de la venta. El margen neto restante será liberado al Comercio según el ciclo de pago elegido (semanal o quincenal).</p>

      <h2 className="text-heading-lg mt-10 mb-4">3. Planes y Facturación</h2>
      <p>La Plataforma se ofrece mediante suscripción mensual o anual (con un descuento del 17%). Los precios vigentes se publican en www.uvico.shop/planes y están expresados en dólares de los Estados Unidos (USD), aunque pueden pagarse en bolívares a la tasa oficial del BCV.</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li><strong>Plan Fundador:</strong> $1/mes por los primeros 3 meses, luego $9/mes. Limitado a los primeros 500 Comercios.</li>
        <li><strong>Plan Esencial:</strong> $19.99/mes.</li>
        <li><strong>Plan Profesional:</strong> $49.99/mes.</li>
        <li><strong>Plan Empresa:</strong> $99.99/mes.</li>
      </ul>
      <p>Los impuestos aplicables (IVA, IGTF) se discriminarán en la factura correspondiente emitida ante el SENIAT.</p>

      <h2 className="text-heading-lg mt-10 mb-4">4. Responsabilidades del Comercio</h2>
      <p>El Comercio es el único vendedor de los bienes o servicios que ofrece en su tienda. Es el único responsable frente a sus Compradores por la calidad, legalidad, disponibilidad, entrega, garantías y facturación de los productos. UVICO no es parte de la compraventa entre el Comercio y el Comprador.</p>
      <p>El Comercio se obliga a suministrar información veraz, incluido su RIF, y a cumplir con las normas de protección al consumidor y tributarias venezolanas.</p>

      <h2 className="text-heading-lg mt-10 mb-4">5. Limitación de responsabilidad</h2>
      <p>UVICO no responde por lucro cesante, pérdida de ventas o daños indirectos derivados del uso de la Plataforma. La responsabilidad total de UVICO se limita al monto pagado por el Comercio en los 3 meses anteriores al hecho que originó el reclamo. UVICO no controla ni responde por las fallas, retenciones o decisiones de aprobación de las Pasarelas de Pago de terceros (Wayu Pay, VeriPay, VEXPay, Payku, Pabilo).</p>
    </article>
  );
}
