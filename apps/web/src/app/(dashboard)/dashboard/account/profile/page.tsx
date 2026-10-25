import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function ProfilePage() {
  return (
    <div className="space-y-6 animate-fade-in max-w-2xl">
      <div>
        <h1 className="text-display-3 font-display text-foreground">Mi Perfil</h1>
        <p className="text-body text-muted-foreground">Edita tu información personal.</p>
      </div>
      <Card className="rounded-2xl border-border shadow-sm">
        <CardHeader><CardTitle>Datos Personales</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Nombre completo</label>
            <input type="text" defaultValue="Yxr" className="w-full h-10 px-3 rounded-lg border border-border bg-white focus:outline-none text-sm" />
          </div>
          <Button className="rounded-full bg-foreground text-background mt-4">Guardar</Button>
        </CardContent>
      </Card>
    </div>
  );
}
