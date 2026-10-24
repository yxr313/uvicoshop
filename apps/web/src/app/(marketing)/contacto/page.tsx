import { Mail, MessageCircle } from "lucide-react";

export default function ContactoPage() {
  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-4xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-16">
        
        <div>
          <h1 className="text-display-2 mb-6 text-foreground tracking-tight">Hablemos.</h1>
          <p className="text-body-lg text-muted-foreground mb-12">
            ¿Tienes alguna duda sobre la plataforma, necesitas soporte técnico o quieres asociarte con nosotros? Estamos aquí para ayudarte.
          </p>

          <div className="space-y-8">
            <div className="flex gap-4">
              <div className="w-10 h-10 bg-neutral-50 border border-border rounded-xl flex items-center justify-center shrink-0">
                <Mail className="w-5 h-5 text-foreground" />
              </div>
              <div>
                <h3 className="text-heading font-display text-foreground mb-1">Escríbenos</h3>
                <p className="text-body-sm text-muted-foreground mb-2">Respondemos en menos de 24 horas hábiles.</p>
                <a href="mailto:contacto@uvico.shop" className="text-body font-medium text-foreground hover:underline">contacto@uvico.shop</a>
              </div>
            </div>

            <div className="flex gap-4">
              <div className="w-10 h-10 bg-neutral-50 border border-border rounded-xl flex items-center justify-center shrink-0">
                <MessageCircle className="w-5 h-5 text-foreground" />
              </div>
              <div>
                <h3 className="text-heading font-display text-foreground mb-1">Comunidad en WhatsApp</h3>
                <p className="text-body-sm text-muted-foreground mb-2">Únete a nuestro canal para novedades.</p>
                <a href="https://whatsapp.com/channel/0029VbEDp1CICVfc8CaOi02w" target="_blank" rel="noreferrer" className="text-body font-medium text-foreground hover:underline">Ir al canal</a>
              </div>
            </div>
          </div>
        </div>

        {/* Formulario de Contacto (Estético) */}
        <div className="bg-neutral-50 p-8 rounded-[2rem] border border-border">
          <form className="space-y-6">
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground">Nombre</label>
              <input type="text" className="w-full h-12 px-4 rounded-xl border border-border bg-white focus:outline-none focus:border-neutral-400 transition-all text-sm" placeholder="Tu nombre" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground">Correo electrónico</label>
              <input type="email" className="w-full h-12 px-4 rounded-xl border border-border bg-white focus:outline-none focus:border-neutral-400 transition-all text-sm" placeholder="tu@correo.com" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground">Mensaje</label>
              <textarea rows={4} className="w-full p-4 rounded-xl border border-border bg-white focus:outline-none focus:border-neutral-400 transition-all text-sm resize-none" placeholder="¿En qué podemos ayudarte?" />
            </div>
            <button type="button" className="w-full bg-foreground text-background h-12 rounded-full font-medium hover:bg-neutral-800 transition-colors">
              Enviar mensaje
            </button>
          </form>
        </div>

      </div>
    </div>
  );
}
