import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function SettingsPage() {
  return (
    <div className="space-y-6 animate-fade-in max-w-2xl">
      <div>
        <h1 className="text-display-3 font-display text-foreground">Ajustes de la Tienda</h1>
        <p className="text-body text-muted-foreground">Configura los detalles de tu negocio.</p>
      </div>
      <Card className="rounded-2xl border-border shadow-sm">
        <CardHeader><CardTitle>Información Básica</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Moneda Base</label>
            <select className="w-full h-10 px-3 rounded-lg border border-border bg-white text-sm">
              <option>USD ($)</option>
              <option>VES (Bs.)</option>
            </select>
          </div>
          <Button className="rounded-full bg-foreground text-background mt-4">Guardar Ajustes</Button>
        </CardContent>
      </Card>
    </div>
  );
}
