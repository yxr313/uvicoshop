import { CheckCircle2 } from "lucide-react";

export default function StatusPage() {
  const services = [
    { name: "API de Uvico", status: "Operativo", uptime: "99.99%" },
    { name: "Panel de Administración", status: "Operativo", uptime: "100%" },
    { name: "Storefront (Tiendas)", status: "Operativo", uptime: "99.98%" },
    { name: "Pasarelas de Pago (Wayu, VeriPay)", status: "Operativo", uptime: "99.95%" },
    { name: "Webhooks y Notificaciones", status: "Operativo", uptime: "100%" },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-4xl mx-auto">
        <div className="mb-16">
          <div className="inline-flex items-center gap-2 bg-green-50 text-green-700 px-4 py-2 rounded-full text-sm font-medium mb-6">
            <CheckCircle2 size={16} /> Todos los sistemas operativos
          </div>
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">Estado del sistema.</h1>
          <p className="text-body-lg text-muted-foreground">Monitoreo en tiempo real de la infraestructura de Uvico.</p>
        </div>

        <div className="border-t border-border">
          {services.map((service, idx) => (
            <div key={idx} className="flex items-center justify-between py-6 border-b border-border">
              <span className="text-heading font-display text-foreground">{service.name}</span>
              <div className="flex items-center gap-8">
                <span className="text-sm text-muted-foreground hidden sm:block">{service.uptime} uptime</span>
                <div className="flex items-center gap-2">
                  <div className="w-2.5 h-2.5 rounded-full bg-green-500"></div>
                  <span className="text-sm font-medium text-foreground">{service.status}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
