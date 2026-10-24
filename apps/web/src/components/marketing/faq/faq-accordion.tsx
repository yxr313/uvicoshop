"use client";

import { useState } from "react";
import { Plus, Minus } from "lucide-react";

export function FAQ() {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const faqs = [
    {
      q: "¿Cómo funciona el dropshipping en Uvico?",
      a: "Tú eliges los productos de nuestra red de proveedores verificados, los publicas en tu tienda con tu propio precio y nosotros gestionamos el pedido con el proveedor cuando un cliente te compre."
    },
    {
      q: "¿Cuándo recibo mi dinero?",
      a: "Retenemos los fondos del dropshipping para protegerte de fraudes y contracargos. Liberamos tus ganancias netas cada viernes directamente a tu cuenta bancaria."
    },
    {
      q: "¿Puedo vender mi propia mercancía?",
      a: "Sí. Si tienes tu propio inventario y te encargas de tus envíos, el dinero de tus ventas va directo a tu pasarela de pago sin retenciones de nuestra parte."
    },
    {
      q: "¿Necesito tarjeta de crédito para la prueba?",
      a: "No. El plan Fundador de 7 días gratis no requiere ningún método de pago. Solo pagas si decides continuar."
    }
  ];

  return (
    <section className="py-32 px-4 bg-background">
      <div className="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-12 gap-12 md:gap-24">
        
        {/* Sidebar Izquierdo (Sticky) */}
        <div className="md:col-span-4">
          <div className="sticky top-32">
            <h2 className="text-display-3 text-foreground">
              Preguntas frecuentes.
            </h2>
          </div>
        </div>

        {/* Acordeón Derecho */}
        <div className="md:col-span-8 space-y-4">
          {faqs.map((faq, index) => {
            const isOpen = openIndex === index;
            return (
              <div 
                key={index} 
                className="border-b border-border pb-4"
              >
                <button
                  onClick={() => setOpenIndex(isOpen ? null : index)}
                  className="w-full py-6 flex items-center justify-between text-left focus:outline-none group"
                >
                  <span className="text-heading font-display text-foreground pr-8 group-hover:text-muted-foreground transition-colors">{faq.q}</span>
                  <span className="text-foreground shrink-0 transition-transform duration-300">
                    {isOpen ? <Minus size={20} /> : <Plus size={20} />}
                  </span>
                </button>
                <div 
                  className={`overflow-hidden transition-all duration-300 ease-in-out ${isOpen ? 'max-h-40 opacity-100' : 'max-h-0 opacity-0'}`}
                >
                  <p className="text-body text-muted-foreground pb-6 pr-12">{faq.a}</p>
                </div>
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
