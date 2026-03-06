import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/lib/auth-context';
import api from '@/lib/api';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Eye, EyeOff, ArrowRight, Loader2 } from 'lucide-react';

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState({
    name: '', email: '', password: '', age: '', emergency_contact: ''
  });

  const handleChange = (e) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const endpoint = isLogin ? '/auth/login' : '/auth/signup';
      const payload = isLogin
        ? { email: form.email, password: form.password }
        : {
            name: form.name,
            email: form.email,
            password: form.password,
            age: parseInt(form.age) || 0,
            emergency_contact: form.emergency_contact
          };

      const response = await api.post(endpoint, payload);
      const { user } = response.data.data;

      // Set auth state (no token needed)
      login(user);

      toast.success(isLogin ? 'Welcome back!' : 'Account created successfully!');
      navigate('/dashboard', { replace: true });
    } catch (error) {
      const msg = error.response?.data?.error || error.response?.data?.message || 'Something went wrong. Please try again.';
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Hero section - desktop only */}
      <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden">
        <img
          src="https://images.unsplash.com/photo-1758518731468-98e90ffd7430?crop=entropy&cs=srgb&fm=jpg&ixid=M3w3NTY2Njl8MHwxfHNlYXJjaHwyfHxkaXZlcnNlJTIwZ3JvdXAlMjBvZiUyMHdvbWVuJTIwbWVldGluZyUyMGNvbmZpZGVudCUyMHByb2Zlc3Npb25hbHxlbnwwfHx8fDE3NzI3MTI3OTd8MA&ixlib=rb-4.1.0&q=85"
          alt="Empowered women professionals"
          className="absolute inset-0 w-full h-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent" />
        <div className="relative z-10 flex flex-col justify-end p-12">
          <h1 className="text-5xl font-bold text-white tracking-tight heading-font">
            SheLeads
          </h1>
          <p className="text-lg text-white/90 mt-3 max-w-md leading-relaxed">
            Empowering women to track their invaluable contributions and grow their businesses with AI-powered tools.
          </p>
        </div>
      </div>

      {/* Form section */}
      <div className="flex-1 flex items-center justify-center p-6 bg-background">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden text-center mb-8">
            <h1 className="text-3xl font-bold text-primary heading-font">SheLeads</h1>
            <p className="text-sm text-muted-foreground mt-1">Empowering Women</p>
          </div>

          <Card className="border-0 shadow-none lg:border lg:shadow-[0_8px_30px_rgb(0,0,0,0.04)]" data-testid="auth-card">
            <CardHeader className="pb-4">
              <CardTitle className="text-2xl heading-font">
                {isLogin ? 'Welcome back' : 'Create account'}
              </CardTitle>
              <CardDescription>
                {isLogin ? 'Sign in to continue your journey' : 'Start empowering your work today'}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                {!isLogin && (
                  <div className="space-y-2">
                    <Label htmlFor="name">Full Name</Label>
                    <Input id="name" name="name" value={form.name} onChange={handleChange} placeholder="Jane Doe" required data-testid="auth-name-input" className="rounded-xl" />
                  </div>
                )}

                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input id="email" name="email" type="email" value={form.email} onChange={handleChange} placeholder="jane@example.com" required data-testid="auth-email-input" className="rounded-xl" />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <div className="relative">
                    <Input
                      id="password" name="password"
                      type={showPassword ? 'text' : 'password'}
                      value={form.password} onChange={handleChange}
                      placeholder="Min. 8 characters" minLength={8} required
                      data-testid="auth-password-input"
                      className="rounded-xl pr-10"
                    />
                    <Button
                      type="button" variant="ghost" size="icon"
                      className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7"
                      onClick={() => setShowPassword(!showPassword)}
                      data-testid="toggle-password-btn"
                    >
                      {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </Button>
                  </div>
                </div>

                {!isLogin && (
                  <>
                    <div className="space-y-2">
                      <Label htmlFor="age">Age</Label>
                      <Input id="age" name="age" type="number" min="1" max="120" value={form.age} onChange={handleChange} placeholder="25" required data-testid="auth-age-input" className="rounded-xl" />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="emergency_contact">Emergency Contact</Label>
                      <Input id="emergency_contact" name="emergency_contact" value={form.emergency_contact} onChange={handleChange} placeholder="+1234567890" required data-testid="auth-emergency-input" className="rounded-xl" />
                    </div>
                  </>
                )}

                <Button type="submit" className="w-full rounded-xl h-11 btn-hover" disabled={loading} data-testid="auth-submit-btn">
                  {loading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                  {isLogin ? 'Sign In' : 'Create Account'}
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>
              </form>

              <div className="mt-6 text-center">
                <button
                  type="button"
                  onClick={() => { setIsLogin(!isLogin); setForm({ name: '', email: '', password: '', age: '', emergency_contact: '' }); }}
                  className="text-sm text-primary hover:underline"
                  data-testid="auth-toggle-btn"
                >
                  {isLogin ? "Don't have an account? Sign up" : "Already have an account? Sign in"}
                </button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
