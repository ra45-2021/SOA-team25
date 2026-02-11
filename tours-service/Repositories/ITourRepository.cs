using TourService.Domain;

namespace TourService.Repositories;

public interface ITourRepository
{
    Task<Tour> AddAsync(Tour tour, CancellationToken ct);
    Task<List<Tour>> GetByAuthorAsync(string authorId, CancellationToken ct);
}
