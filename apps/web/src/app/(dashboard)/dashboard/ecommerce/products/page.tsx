"use client";

import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Search, Plus, MoreVertical, SlidersHorizontal, Loader2 } from "lucide-react";

type Product = {
  id: string;
  name: string;
  description: string;
  price_usd: number;
};

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  // Fetch a la API de Go
  const fetchProducts = async () => {
    try {
      setLoading(true);
      const res = await fetch("http://localhost:8080/api/v1/products");
      const json = await res.json();
      setProducts(json.data || []);
    } catch (error) {
      console.error("Error cargando productos", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-display-3 font-display text-foreground">Productos</h1>
          <p className="text-body text-muted-foreground">Gestiona tu catálogo real.</p>
        </div>
        <Button className="rounded-full bg-accent text-accent-foreground hover:bg-accent-hover">
          <Plus size={16} className="mr-2" /> Añadir Producto
        </Button>
      </div>

      <Card className="rounded-2xl border-border shadow-sm overflow-hidden">
        <div className="p-4 border-b border-border flex justify-between items-center bg-white">
          <div className="relative w-72">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
            <input type="text" placeholder="Buscar productos..." className="w-full h-9 pl-9 pr-4 rounded-md border border-border bg-transparent text-sm focus:outline-none" />
          </div>
        </div>
        <CardContent className="p-0">
          {loading ? (
            <div className="h-[300px] flex items-center justify-center">
              <Loader2 className="animate-spin text-muted-foreground" size={32} />
            </div>
          ) : products.length === 0 ? (
            <div className="h-[300px] flex items-center justify-center text-muted-foreground">
              No hay productos. ¡Crea el primero!
            </div>
          ) : (
            <table className="w-full text-left text-sm">
              <thead className="bg-neutral-50 border-b border-border">
                <tr>
                  <th className="px-6 py-3 font-medium text-muted-foreground">Producto</th>
                  <th className="px-6 py-3 font-medium text-muted-foreground">Precio</th>
                  <th className="px-6 py-3 font-medium text-muted-foreground">Estado</th>
                  <th className="px-6 py-3"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {products.map((p) => (
                  <tr key={p.id} className="hover:bg-neutral-50/50">
                    <td className="px-6 py-4 flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-neutral-100" />
                      <div>
                        <p className="font-medium text-foreground">{p.name}</p>
                        <p className="text-xs text-muted-foreground">{p.description}</p>
                      </div>
                    </td>
                    <td className="px-6 py-4 font-medium text-foreground">${p.price_usd}</td>
                    <td className="px-6 py-4">
                      <span className="px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-700">Activo</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <MoreVertical size={16} className="text-muted-foreground cursor-pointer" />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
