import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function PricingPage() {
  return (
    <div className="space-y-6 animate-fade-in">
      <div>
        <h1 className="text-display-3 font-display text-foreground">Planes y Facturación</h1>
        <p className="text-body text-muted-foreground">Gestiona tu suscripción.</p>
      </div>
      <Card className="rounded-2xl border-foreground border-2 shadow-lg max-w-md">
        <CardHeader><CardTitle>Plan Fundador</CardTitle></CardHeader>
        <CardContent>
          <p className="text-3xl font-display font-bold mb-4">$9<span className="text-sm text-muted-foreground font-normal">/mes</span></p>
          <Button className="w-full rounded-full bg-accent text-accent-foreground">Gestionar Suscripción</Button>
        </CardContent>
      </Card>
    </div>
  );
}
