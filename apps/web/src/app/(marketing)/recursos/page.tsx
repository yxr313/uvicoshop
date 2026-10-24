import Link from "next/link";
import { ArrowRight, PlayCircle, BookOpen, GraduationCap, Youtube } from "lucide-react";

export default function RecursosPage() {
  const resources = [
    {
      category: "Guía esencial",
      title: "Cómo abrir tu primera tienda en 10 minutos",
      desc: "El paso a paso definitivo para configurar tu catálogo, conectar pagos y hacer tu primera venta.",
      icon: <BookOpen className="w-5 h-5" />,
      videoUrl: "https://www.youtube.com/results?search_query=c%C3%B3mo+crear+tienda+online+dropshipping+venezuela",
    },
    {
      category: "Video tutorial",
      title: "Dominando el editor visual Puck",
      desc: "Aprende a crear landing pages de productos de alta conversión usando bloques modulares.",
      icon: <PlayCircle className="w-5 h-5" />,
      videoUrl: "https://www.youtube.com/results?search_query=editor+visual+puck+tutorial+espa%C3%B1ol",
    },
    {
      category: "Estrategia",
      title: "Cómo elegir productos ganadores",
      desc: "Los criterios exactos que usan los grandes dropshippers para validar un producto antes de invertir.",
      icon: <GraduationCap className="w-5 h-5" />,
      videoUrl: "https://www.youtube.com/results?search_query=productos+ganadores+dropshipping+espa%C3%B1ol",
    },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-7xl mx-auto">
        <div className="max-w-3xl mb-20">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">Aprende a escalar.</h1>
          <p className="text-body-lg text-muted-foreground">Guías, tutoriales y estrategias para llevar tu e-commerce al siguiente nivel.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {resources.map((post, idx) => (
            <div key={idx} className="flex flex-col h-full">
              <div className="flex-1 border-t border-border pt-6 mb-6">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
                  {post.icon} {post.category}
                </div>
                <h3 className="text-heading-lg font-display text-foreground mb-3">{post.title}</h3>
                <p className="text-body-sm text-muted-foreground mb-6">{post.desc}</p>
              </div>
              <a href={post.videoUrl} target="_blank" rel="noopener noreferrer" className="mt-auto flex items-center gap-2 text-sm font-medium text-foreground hover:text-accent transition-colors group">
                <Youtube className="w-5 h-5" /> Ver video tutorial <ArrowRight className="w-4 h-4 ml-1 group-hover:translate-x-1 transition-transform" />
              </a>
            </div>
          ))}
        </div>

        <div className="mt-32 bg-neutral-50 rounded-[2rem] p-12 text-center border border-border">
          <h2 className="text-display-3 text-foreground mb-4">¿Buscas ayuda técnica?</h2>
          <p className="text-body text-muted-foreground mb-8 max-w-md mx-auto">Visita nuestro centro de ayuda o ponte en contacto directo con el equipo de soporte.</p>
          <a href="mailto:contacto@uvico.shop"><button className="bg-foreground text-background px-8 py-3 rounded-full font-medium hover:bg-neutral-800 transition-colors">Contactar a soporte</button></a>
        </div>
      </div>
    </div>
  );
}
