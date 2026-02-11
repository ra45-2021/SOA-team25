using Microsoft.EntityFrameworkCore;
using TourService.Data;
using TourService.Domain;

namespace TourService.Repositories;

public class TourRepository : ITourRepository
{
    private readonly ToursContext _db;
    public TourRepository(ToursContext db) { _db = db; }

    public async Task<Tour> AddAsync(Tour tour, CancellationToken ct)
    {
        _db.Tours.Add(tour);
        await _db.SaveChangesAsync(ct);
        return tour;
    }

    public Task<List<Tour>> GetByAuthorAsync(string authorId, CancellationToken ct)
    {
        return _db.Tours
            .AsNoTracking()
            .Include(t => t.TourTags)
                .ThenInclude(tt => tt.Tag)
            .Where(t => t.AuthorId == authorId)
            .OrderByDescending(t => t.CreatedAtUtc)
            .ToListAsync(ct);
    }
}
