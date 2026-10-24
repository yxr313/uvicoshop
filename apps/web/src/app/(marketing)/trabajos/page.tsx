import Link from "next/link";
import { ArrowRight, MapPin, Clock } from "lucide-react";

export default function TrabajosPage() {
  const jobs = [
    {
      title: "Senior Backend Engineer (Go)",
      dept: "Ingeniería",
      location: "Remoto (LATAM)",
      type: "Tiempo completo",
    },
    {
      title: "Frontend Engineer (Next.js)",
      dept: "Ingeniería",
      location: "Remoto (LATAM)",
      type: "Tiempo completo",
    },
    {
      title: "Product Designer",
      dept: "Diseño",
      location: "Remoto (LATAM)",
      type: "Tiempo completo",
    },
    {
      title: "Growth Marketing Manager",
      dept: "Marketing",
      location: "Caracas, Venezuela / Remoto",
      type: "Tiempo completo",
    },
    {
      title: "Customer Success Specialist",
      dept: "Soporte",
      location: "Remoto (LATAM)",
      type: "Tiempo completo",
    },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-5xl mx-auto">
        
        {/* Header */}
        <div className="max-w-3xl mb-20">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">
            Construye el futuro del e-commerce.
          </h1>
          <p className="text-body-lg text-muted-foreground">
            Estamos buscando personas obsesionadas con el detalle, que quieran resolver problemas complejos para impulsar a miles de emprendedores en Venezuela y LATAM.
          </p>
        </div>

        {/* Lista de Empleos */}
        <div className="mb-24">
          <h2 className="text-heading-lg font-display text-foreground mb-8">Posiciones abiertas</h2>
          <div className="border-t border-border">
            {jobs.map((job, idx) => (
              <a 
                key={idx} 
                href={`mailto:yxr@uvico.shop?subject=Aplicación para ${job.title}`}
                className="group flex flex-col md:flex-row md:items-center justify-between py-6 border-b border-border hover:bg-neutral-50/50 transition-colors px-4 -mx-4 rounded-lg"
              >
                <div>
                  <h3 className="text-heading font-display text-foreground mb-2 group-hover:text-accent transition-colors">
                    {job.title}
                  </h3>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    <span className="font-medium text-foreground">{job.dept}</span>
                    <span className="flex items-center gap-1"><MapPin size={14} /> {job.location}</span>
                    <span className="flex items-center gap-1 hidden sm:flex"><Clock size={14} /> {job.type}</span>
                  </div>
                </div>
                <div className="mt-4 md:mt-0 flex items-center text-sm font-medium text-foreground">
                  Aplicar ahora <ArrowRight className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
                </div>
              </a>
            ))}
          </div>
        </div>

        {/* Cultura / Valores */}
        <div className="border-t border-border pt-20">
          <h2 className="text-display-2 text-foreground mb-12">Nuestra cultura.</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            <div>
              <h3 className="text-heading font-display text-foreground mb-3">Propiedad radical</h3>
              <p className="text-body-sm text-muted-foreground">No esperamos órdenes. Detectamos problemas, proponemos soluciones y las ejecutamos hasta el final.</p>
            </div>
            <div>
              <h3 className="text-heading font-display text-foreground mb-3">Obsesión por el usuario</h3>
              <p className="text-body-sm text-muted-foreground">Cada decisión de diseño o de código debe responder a una necesidad real del emprendedor venezolano.</p>
            </div>
            <div>
              <h3 className="text-heading font-display text-foreground mb-3">Velocidad con calidad</h3>
              <p className="text-body-sm text-muted-foreground">Iteramos rápido, pero con estándares de Silicon Valley. El código que escribimos hoy no nos dará vergüenza mañana.</p>
            </div>
          </div>
        </div>

      </div>
    </div>
  );
}
