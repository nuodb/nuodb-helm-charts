
select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;


select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

select /* goals scored in single year for Gordie Howe, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'howego01'
and   s.year     = 1946;

select /* goals scored in single year for Wayne Gretzky, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'gretzwa01'
and   s.goals    > 10
and   s.year     = 1978;

select /* goals scored in single year for Bobbie Hull, first row */ 
  firstname, lastname, birthyear, s.year "YEAR PLAYED", t.name "TEAM", s.goals "GOALS"
from scoring s, players p, teams t
where s.playerid = p.playerid
and   s.year     = t.year
and   s.teamid   = t.teamid
and   s.playerid = 'hullbo01'
and   s.year     = 1957;

update scoring
set goals = goals
where playerid = 'hullbo01'
and   year = 1957;

