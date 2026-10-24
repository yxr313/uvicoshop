import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ShieldCheck, Banknote, Clock, Truck, Search, CheckCircle, HelpCircle } from "lucide-react";

export default function ProveedoresPage() {
  const verificationSteps = [
    { icon: <Search className="w-6 h-6" />, title: "1. Auditoría", desc: "Revisamos el historial, la capacidad de inventario y la reputación del proveedor." },
    { icon: <CheckCircle className="w-6 h-6" />, title: "2. Prueba de calidad", desc: "Hacemos pedidos de prueba para verificar tiempos de despacho y embalaje." },
    { icon: <ShieldCheck className="w-6 h-6" />, title: "3. Contrato", desc: "Firmamos acuerdos de nivel de servicio (SLA) para garantizar entregas a tiempo." },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-7xl mx-auto">
        <div className="max-w-3xl mb-20">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">Dropshipping sin riesgos.</h1>
          <p className="text-body-lg text-muted-foreground">Nosotros verificamos a los proveedores. Tú te enfocas en vender.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-16 items-center mb-24">
          <div>
            <h2 className="text-display-2 text-foreground mb-6">Red de proveedores verificados.</h2>
            <p className="text-body-lg text-muted-foreground mb-8">No tienes que buscar proveedores en Instagram ni arriesgarte a estafas. Accede directamente a nuestra red de mayoristas y fabricantes nacionales con inventario real y despacho garantizado.</p>
            <div className="space-y-6">
              <div className="flex gap-4"><ShieldCheck className="w-6 h-6 text-foreground shrink-0" /><div><h4 className="text-body font-medium text-foreground">Verificación exhaustiva</h4><p className="text-sm text-muted-foreground">Auditamos stock y tiempos de entrega antes de integrarlos.</p></div></div>
              <div className="flex gap-4"><Truck className="w-6 h-6 text-foreground shrink-0" /><div><h4 className="text-body font-medium text-foreground">Logística integrada</h4><p className="text-sm text-muted-foreground">Conexión con MRW, Zoom y delivery propio.</p></div></div>
            </div>
          </div>
          <div className="bg-neutral-50 rounded-[2rem] p-8 border border-border">
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <div key={i} className="bg-white p-4 rounded-xl border border-border flex items-center justify-between shadow-sm">
                  <div className="flex items-center gap-4"><div className="w-10 h-10 bg-neutral-100 rounded-lg" /><div><div className="w-24 h-3 bg-neutral-200 rounded mb-2" /><div className="w-16 h-2 bg-neutral-100 rounded" /></div></div>
                  <div className="bg-green-50 text-green-700 text-xs font-semibold px-2 py-1 rounded">Verificado</div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="border-t border-border pt-20 mb-24">
          <h2 className="text-display-2 text-foreground mb-12">Tu dinero, protegido.</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            <div><Banknote className="w-8 h-8 text-foreground mb-4" /><h3 className="text-heading font-display text-foreground mb-2">1. El cliente paga</h3><p className="text-body text-muted-foreground">El dinero cae en una cuenta segura de Uvico, no en la del proveedor.</p></div>
            <div><ShieldCheck className="w-8 h-8 text-foreground mb-4" /><h3 className="text-heading font-display text-foreground mb-2">2. Protegemos los fondos</h3><p className="text-body text-muted-foreground">Retenemos el 100% del pago para evitar fraudes o contracargos.</p></div>
            <div><Clock className="w-8 h-8 text-foreground mb-4" /><h3 className="text-heading font-display text-foreground mb-2">3. Liberamos tu ganancia</h3><p className="text-body text-muted-foreground">Cada viernes te transferimos tu margen neto (menos el 5% de comisión) a tu banco.</p></div>
          </div>
        </div>

        <div className="border-t border-border pt-20 mb-24">
          <h2 className="text-display-2 mb-12">Proceso de verificación.</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            {verificationSteps.map((step, idx) => (
              <div key={idx}><div className="w-12 h-12 bg-neutral-50 rounded-full flex items-center justify-center mb-4">{step.icon}</div><h3 className="text-heading font-display text-foreground mb-2">{step.title}</h3><p className="text-body-sm text-muted-foreground">{step.desc}</p></div>
            ))}
          </div>
        </div>

        <div className="border-t border-border pt-20">
          <h2 className="text-display-2 mb-12">Preguntas frecuentes sobre proveedores.</h2>
          <div className="space-y-6 max-w-3xl">
            <div><h4 className="text-heading font-display text-foreground mb-2">¿Puedo elegir mis propios proveedores?</h4><p className="text-body text-muted-foreground">Sí. Si ya tienes un proveedor de confianza, puedes agregarlo manualmente. Sin embargo, no estará cubierto por nuestra garantía de custodia de fondos.</p></div>
            <div><h4 className="text-heading font-display text-foreground mb-2">¿Qué pasa si un proveedor no despacha a tiempo?</h4><p className="text-body text-muted-foreground">Si el proveedor incumple el SLA, retenemos el pago y te notificamos de inmediato. Puedes cancelar el pedido sin penalización.</p></div>
            <div><h4 className="text-heading font-display text-foreground mb-2">¿Los proveedores envían a todo el país?</h4><p className="text-body text-muted-foreground">La mayoría de nuestra red cubre las principales ciudades. Puedes filtrar proveedores por zona de cobertura.</p></div>
          </div>
        </div>
      </div>
    </div>
  );
}
