using TourService.Contracts;

namespace TourService.Services;

public interface ITourService
{
    Task<TourDto> CreateAsync(string authorId, CreateTourRequest req, CancellationToken ct);
    Task<List<TourDto>> GetMyToursAsync(string authorId, CancellationToken ct);
}
