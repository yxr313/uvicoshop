export default function ReembolsosPage() {
  return (
    <article>
      <h1 className="text-display-3 mb-2">Política de Reembolsos y Cancelación</h1>
      <p className="text-sm text-muted-foreground mb-12">Versión 1.1 | Vigente desde la fecha de publicación</p>

      <h2 className="text-heading-lg mt-10 mb-4">1. Alcance</h2>
      <p>Esta Política regula la cancelación de suscripciones a UVICO y los reembolsos de los cobros a los Comercios. <strong>No regula las compras que un Comprador hace en la tienda de un Comercio</strong>; esas devoluciones son responsabilidad exclusiva del Comercio vendedor, bajo las reglas detalladas en la sección 2.</p>

      <h2 className="text-heading-lg mt-10 mb-4">2. Reembolsos a Compradores (Compras)</h2>
      <p>Para proteger a ambas partes y cumplir con las normativas de protección al consumidor, aplicamos una política de dos niveles:</p>
      <ul className="list-disc pl-5 space-y-2 mb-6">
        <li><strong>Cambio de opinión:</strong> Si el Comprador decide no quedarse con el producto, no se devuelve dinero en efectivo. Se otorgará un cupón por el 100% del monto, válido para cualquier tienda dentro de UVICO. El costo es absorbido por el saldo retenido del Comercio.</li>
        <li><strong>Culpa del Comercio (No entrega, defectuoso o distinto):</strong> El reembolso en dinero es obligatorio y se debitará automáticamente del saldo retenido del Comercio. Esta medida previene sanciones de SUNDDE y contracargos bancarios.</li>
        <li><strong>Contracargos (Chargebacks):</strong> Si el Comprador abre una disputa en su banco, UVICO debitará el monto del saldo del Comercio. Si el saldo no alcanza, se descontará de futuras liberaciones de fondos.</li>
      </ul>

      <h2 className="text-heading-lg mt-10 mb-4">3. Cancelación de Suscripción UVICO</h2>
      <p>El Comercio puede cancelar su suscripción en cualquier momento desde el panel de control, sin penalidad. La cancelación detiene los cobros futuros y la tienda sigue activa hasta el final del Período ya pagado.</p>
      <p><strong>No se realizan reembolsos proporcionales</strong> por el tiempo no usado de un mes o año pagado, ni se reembolsa la comisión del 5% por ventas ya procesadas.</p>

      <h2 className="text-heading-lg mt-10 mb-4">4. Reembolsos por errores de UVICO</h2>
      <p>UVICO reembolsará el cobro cuando: (a) se haya cobrado doble la misma cuota; (b) el monto cobrado difiera del plan contratado; o (c) se haya cobrado un período posterior a una cancelación solicitada a tiempo. Las solicitudes deben enviarse a <strong>soporte@uvico.shop</strong> dentro de los 15 días siguientes al cobro.</p>
    </article>
  );
}
