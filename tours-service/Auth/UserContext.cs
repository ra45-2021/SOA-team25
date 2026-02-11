using System.Security.Claims;

namespace TourService.Auth;

public static class UserContext
{
    public static string GetUserId(ClaimsPrincipal user)
    {
        return user.FindFirstValue(ClaimTypes.NameIdentifier)
            ?? user.FindFirstValue("sub")
            ?? "dev-user";
    }
}
