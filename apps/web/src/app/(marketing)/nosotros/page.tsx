import Image from "next/image";

export default function NosotrosPage() {
  const values = [
    { title: "Obsesión por el detalle", desc: "Creemos que el diseño y el código limpio no son un lujo, son la base de la confianza. Si el botón no se siente bien, no está terminado." },
    { title: "Velocidad brutal", desc: "Lanzamos, medimos y aprendemos. No esperamos seis meses para lanzar una función perfecta; iteramos rápido para resolver problemas reales hoy." },
    { title: "Soberanía tecnológica", desc: "Venezuela no puede depender de plataformas extranjeras que no entienden nuestra economía. Construimos infraestructura local adaptada a nuestra realidad." },
    { title: "Propiedad radical", desc: "Cada miembro del equipo actúa como dueño. No hay espacio para el 'no es mi trabajo'. Si algo está roto, lo arreglamos." }
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-6xl mx-auto">
        
        {/* Hero Statement */}
        <div className="max-w-4xl mb-32">
          <h1 className="text-display-1 mb-8 text-foreground tracking-tight">
            Estamos construyendo el sistema operativo del comercio venezolano.
          </h1>
          <p className="text-[1.5rem] leading-snug text-muted-foreground font-light">
            Uvico no es solo una herramienta para crear tiendas. Es la infraestructura financiera, logística y tecnológica que los emprendedores de LATAM necesitan para competir en la economía global.
          </p>
        </div>

        {/* Cita Filosófica */}
        <div className="mb-32 border-y border-border py-16">
          <blockquote className="max-w-4xl mx-auto text-center">
            <p className="text-display-3 font-display text-foreground italic mb-6">
              "El desarrollo no puede medirse por la cantidad de cosas que se producen, sino por la libertad de las personas para crear su propio destino."
            </p>
            <footer className="text-body font-medium text-muted-foreground tracking-widest uppercase">
              — Ivan Illich
            </footer>
          </blockquote>
        </div>

        {/* La Historia */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-16 mb-32">
          <div>
            <h2 className="text-display-3 text-foreground mb-6">El problema</h2>
            <p className="text-body-lg text-muted-foreground">
              Hacer dropshipping o vender online en Venezuela históricamente ha sido un caos. Docenas de herramientas fragmentadas, pasarelas de pago que fallan, proveedores estafadores y plataformas en inglés que no entienden lo que es un Pago Móvil o una tasa del BCV.
            </p>
          </div>
          <div>
            <h2 className="text-display-3 text-foreground mb-6">La solución</h2>
            <p className="text-body-lg text-muted-foreground">
              Decidimos construir el ecosistema completo en un solo lugar. Un editor visual de clase mundial, un sistema de custodia de fondos que protege al comerciante (y al comprador), y una red de proveedores auditados. Todo diseñado con los estándares de Silicon Valley, pero con el ADN de Venezuela.
            </p>
          </div>
        </div>

        {/* Valores */}
        <div className="mb-32">
          <h2 className="text-display-2 text-foreground mb-12">Nuestros principios.</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-16 gap-y-12">
            {values.map((val, idx) => (
              <div key={idx} className="border-t border-border pt-6">
                <h3 className="text-heading-lg font-display text-foreground mb-3">{val.title}</h3>
                <p className="text-body text-muted-foreground">{val.desc}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Fundador con Foto y Nueva Nota */}
        <div className="border-t border-border pt-20">
          <div className="flex flex-col md:flex-row items-start md:items-center gap-12">
            
            {/* Foto del Fundador */}
            <div className="relative w-48 h-48 shrink-0 rounded-full overflow-hidden border border-border bg-neutral-100">
              <Image 
                src="/foto.jpg" 
                alt="Yxr - Fundador de Uvico" 
                fill
                className="object-cover"
                priority
              />
            </div>

            <div>
              <h2 className="text-display-3 font-display text-foreground mb-4">Nota del fundador</h2>
              <p className="text-body-lg text-muted-foreground max-w-3xl mb-6">
                "No creo en la suerte ni en depender de plataformas extranjeras que ignoran nuestra realidad. Creo en el trabajo, en el código limpio y en la soberanía económica de Venezuela. Uvico no es un experimento; es la base sobre la que se construirá la próxima generación de negocios en LATAM. Tenemos el talento, ahora tenemos las herramientas."
              </p>
              <div>
                <p className="text-body font-medium text-foreground">Yxr</p>
                <p className="text-sm text-muted-foreground">Founder & CEO</p>
              </div>
            </div>

          </div>
        </div>

      </div>
    </div>
  );
}
