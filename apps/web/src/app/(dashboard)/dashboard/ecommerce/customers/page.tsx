import { Card, CardContent } from "@/components/ui/card";

export default function CustomersPage() {
  return (
    <div className="space-y-6 animate-fade-in">
      <div>
        <h1 className="text-display-3 font-display text-foreground">Clientes</h1>
        <p className="text-body text-muted-foreground">Base de datos de tus compradores.</p>
      </div>
      <Card className="rounded-2xl border-border shadow-sm">
        <CardContent className="h-[400px] flex items-center justify-center">
          <p className="text-muted-foreground text-sm">Aquí irá la tabla de Clientes.</p>
        </CardContent>
      </Card>
    </div>
  );
}
