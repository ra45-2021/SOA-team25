using Ocelot.DependencyInjection;
using Ocelot.Middleware;
using ApiGateway.Protos; 
using Microsoft.AspNetCore.Mvc; 

var builder = WebApplication.CreateBuilder(args);

builder.Configuration.AddJsonFile("ocelot.json", optional: false, reloadOnChange: true);
builder.Services.AddOcelot();

builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(p => p.AllowAnyOrigin().AllowAnyHeader().AllowAnyMethod());
});

builder.Services.AddGrpcClient<AuthService.AuthServiceClient>(o =>
{
    o.Address = new Uri("http://auth-service:50051");
});

var app = builder.Build();

app.UseCors();

app.MapGet("/health", () => "ok");

app.MapGet("/api/profiles/{id}", async (long id, [FromServices] AuthService.AuthServiceClient client) =>
{
    try
    {
        var reply = await client.GetUserAsync(new GetUserRequest { Id = id });
        return Results.Ok(new { 
            id = reply.Id, 
            username = reply.Username, 
            role = reply.Role 
        });
    }
    catch (Exception ex)
    {
        return Results.Problem(ex.Message);
    }
});

await app.UseOcelot();

app.Run();