using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using TourService.Auth;
using TourService.Contracts;
using TourService.Services;

namespace TourService.Controllers;

[ApiController]
[Route("api/tours")]
public class ToursController : ControllerBase
{
    private readonly ITourService _service;

    public ToursController(ITourService service)
    {
        _service = service;
    }

    [HttpPost]
    [Authorize]
    public async Task<ActionResult<TourDto>> Create([FromBody] CreateTourRequest req, CancellationToken ct)
    {
        var authorId = UserContext.GetUserId(User);
        var created = await _service.CreateAsync(authorId, req, ct);
        return CreatedAtAction(nameof(GetMine), new { }, created);
    }

    [HttpGet("mine")]
    [Authorize]
    public async Task<ActionResult<List<TourDto>>> GetMine(CancellationToken ct)
    {
        var authorId = UserContext.GetUserId(User);
        var tours = await _service.GetMyToursAsync(authorId, ct);
        return Ok(tours);
    }
}
