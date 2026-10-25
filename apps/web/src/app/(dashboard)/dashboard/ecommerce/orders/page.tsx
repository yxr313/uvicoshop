import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Search, Plus } from "lucide-react";

export default function OrdersPage() {
  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-display-3 font-display text-foreground">Pedidos</h1>
          <p className="text-body text-muted-foreground">Gestiona los pedidos de tu tienda.</p>
        </div>
        <Button className="rounded-full bg-accent text-accent-foreground hover:bg-accent-hover">
          <Plus size={16} className="mr-2" /> Crear Pedido
        </Button>
      </div>
      <Card className="rounded-2xl border-border shadow-sm overflow-hidden">
        <div className="p-4 border-b border-border flex items-center bg-white">
          <div className="relative w-72">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
            <input type="text" placeholder="Buscar pedidos..." className="w-full h-9 pl-9 pr-4 rounded-md border border-border bg-transparent text-sm focus:outline-none" />
          </div>
        </div>
        <CardContent className="h-[400px] flex items-center justify-center">
          <p className="text-muted-foreground text-sm">Aquí irá la tabla de Pedidos.</p>
        </CardContent>
      </Card>
    </div>
  );
}
