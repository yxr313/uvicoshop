import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function SecurityPage() {
  return (
    <div className="space-y-6 animate-fade-in max-w-2xl">
      <div>
        <h1 className="text-display-3 font-display text-foreground">Seguridad</h1>
        <p className="text-body text-muted-foreground">Protege tu cuenta.</p>
      </div>
      <Card className="rounded-2xl border-border shadow-sm">
        <CardHeader><CardTitle>Cambiar Contraseña</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Nueva contraseña</label>
            <input type="password" placeholder="••••••••" className="w-full h-10 px-3 rounded-lg border border-border bg-white focus:outline-none text-sm" />
          </div>
          <Button className="rounded-full bg-foreground text-background mt-4">Actualizar</Button>
        </CardContent>
      </Card>
    </div>
  );
}
