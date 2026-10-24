import Link from "next/link";
import { ArrowRight } from "lucide-react";

export default function BlogPage() {
  const posts = [
    {
      category: "Estrategia",
      title: "Cómo validar un producto antes de invertir un solo dólar",
      excerpt: "El framework exacto que usan los grandes dropshippers para saber si un producto va a vender o si es una pérdida de tiempo.",
      date: "24 Oct, 2026",
      readTime: "5 min de lectura",
    },
    {
      category: "Finanzas",
      title: "Cómo manejar la tasa del BCV en tu tienda sin perder margen",
      excerpt: "La guía definitiva para fijar precios en dólares, cobrar en bolívares y no morir en el intento con la inflación.",
      date: "20 Oct, 2026",
      readTime: "7 min de lectura",
    },
    {
      category: "Diseño",
      title: "Por qué tu tienda necesita menos botones y más espacio en blanco",
      excerpt: "Análisis de conversión (CRO): cómo el minimalismo y la limpieza visual aumentan la confianza del comprador.",
      date: "15 Oct, 2026",
      readTime: "4 min de lectura",
    },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-5xl mx-auto">
        
        <div className="max-w-3xl mb-20">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">El Blog de Uvico.</h1>
          <p className="text-body-lg text-muted-foreground">
            Estrategias de e-commerce, finanzas y diseño para escalar tu negocio en Venezuela.
          </p>
        </div>

        <div className="border-t border-border">
          {posts.map((post, idx) => (
            <Link 
              key={idx} 
              href="#" 
              className="group flex flex-col py-10 border-b border-border hover:bg-neutral-50/50 transition-colors px-4 -mx-4 rounded-xl"
            >
              <div className="flex items-center gap-3 text-sm font-medium text-muted-foreground mb-4">
                <span className="text-foreground">{post.category}</span>
                <span>•</span>
                <span>{post.date}</span>
                <span>•</span>
                <span>{post.readTime}</span>
              </div>
              <h2 className="text-display-3 font-display text-foreground mb-4 group-hover:text-accent transition-colors">
                {post.title}
              </h2>
              <p className="text-body-lg text-muted-foreground max-w-3xl mb-6">
                {post.excerpt}
              </p>
              <div className="flex items-center text-sm font-medium text-foreground mt-auto">
                Leer artículo <ArrowRight className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
              </div>
            </Link>
          ))}
        </div>

      </div>
    </div>
  );
}
